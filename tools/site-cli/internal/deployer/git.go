package deployer

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DeployOptions 배포 실행 옵션
type DeployOptions struct {
	DistDir string // 배포할 정적 산출물 디렉터리 (기본: "dist")
	Branch  string // 대상 GitHub Pages 브랜치 (기본: "gh-pages")
	Remote  string // Git 리모트 명칭 (기본: "origin")
	CNAME   string // 자동 생성할 CNAME 도메인 (기본: "www.joinc.co.kr")
	Message string // 커밋 메시지
}

// Deploy 정적 사이트를 GitHub Pages 브랜치로 격리 푸시합니다.
func Deploy(opts DeployOptions) error {
	if opts.DistDir == "" {
		opts.DistDir = "dist"
	}
	if opts.Branch == "" {
		opts.Branch = "gh-pages"
	}
	if opts.Remote == "" {
		opts.Remote = "origin"
	}
	if opts.Message == "" {
		opts.Message = "deploy: publish static site via site-cli"
	}

	// 1. dist 디렉터리 존재 여부 확인
	if info, err := os.Stat(opts.DistDir); err != nil || !info.IsDir() {
		return fmt.Errorf("deployer: output directory %q does not exist. Run 'site-cli build' first", opts.DistDir)
	}

	// 2. CNAME 파일 생성
	if opts.CNAME != "" {
		cnamePath := filepath.Join(opts.DistDir, "CNAME")
		if err := os.WriteFile(cnamePath, []byte(strings.TrimSpace(opts.CNAME)+"\n"), 0644); err != nil {
			return fmt.Errorf("deployer: failed to write CNAME file: %w", err)
		}
		log.Printf("[DEPLOY] Injected CNAME: %s", opts.CNAME)
	}

	// 3. 부모 리포지토리의 원격 저장소 URL 획득
	remoteURLCmd := exec.Command("git", "config", "--get", fmt.Sprintf("remote.%s.url", opts.Remote))
	urlBytes, err := remoteURLCmd.Output()
	if err != nil {
		return fmt.Errorf("deployer: failed to get git remote url for %q: %w", opts.Remote, err)
	}
	remoteURL := strings.TrimSpace(string(urlBytes))
	if remoteURL == "" {
		return fmt.Errorf("deployer: remote url for %q is empty", opts.Remote)
	}

	// 4. dist 디렉터리 내부에 독립 임시 Git 저장소를 구성하여 Orphan Push 수행
	// [P-06 & Clean GitOps] 로컬 작업 트리에 일체 간섭하지 않는 완전 격리 배포
	distGitDir := filepath.Join(opts.DistDir, ".git")
	_ = os.RemoveAll(distGitDir) // 이전 임시 .git 정리

	commands := [][]string{
		{"git", "init"},
		{"git", "checkout", "-b", opts.Branch},
		{"git", "add", "-A"},
		{"git", "commit", "-m", opts.Message},
		{"git", "remote", "add", "target", remoteURL},
		{"git", "push", "--force", "target", opts.Branch},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = opts.DistDir
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("deployer: command failed %q: %w (output: %s)", strings.Join(cmdArgs, " "), err, string(out))
		}
	}

	// 배포 완료 후 dist/.git 정리
	_ = os.RemoveAll(distGitDir)

	log.Printf("[SUCCESS] Successfully deployed %s to %s/%s", opts.DistDir, remoteURL, opts.Branch)
	if opts.CNAME != "" {
		log.Printf("[SUCCESS] Custom domain live at: https://%s", opts.CNAME)
	}

	return nil
}
