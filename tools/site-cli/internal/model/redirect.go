package model

// QueryRedirect 쿼리 파라미터 기반 리다이렉트 설정
type QueryRedirect struct {
	SourcePath    string            `yaml:"source_path"`
	QueryKey      string            `yaml:"query_key"`
	DefaultTarget string            `yaml:"default_target"`
	Mappings      map[string]string `yaml:"mappings"`
}

// PathRedirect 정적 경로 기반 1:1 리다이렉트 설정
type PathRedirect struct {
	SourcePath string `yaml:"source_path"`
	TargetPath string `yaml:"target_path"`
}

// RedirectConfig redirect.yaml 전체 설정 구조체
type RedirectConfig struct {
	QueryRedirects []QueryRedirect `yaml:"query_redirects"`
	PathRedirects  []PathRedirect  `yaml:"path_redirects"`
}
