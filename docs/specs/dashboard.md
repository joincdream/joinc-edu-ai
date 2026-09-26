# 📝 포스트 배포 관리 대시보드

```dataview
TABLE 
    created_date AS "생성일", 
    published_date AS "배포일", 
    publish_link AS "배포 링크", 
    tags AS "태그",
    post_id AS "포스트 ID"
FROM "posts/drafts"
SORT created_date DESC
```
