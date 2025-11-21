package test

// Example usage of the new query filtering functionality
// This demonstrates how to use the enhanced GetPosts endpoint

/*
Example API calls:

1. Basic pagination:
GET /api/v1/posts?limit=10&offset=0

2. Search functionality:
GET /api/v1/posts?search=golang&limit=20

3. Sorting by different fields:
GET /api/v1/posts?sort_by=title&sort_order=asc&limit=10
GET /api/v1/posts?sort_by=created_at&sort_order=desc&limit=10
GET /api/v1/posts?sort_by=view_count&sort_order=desc&limit=10

4. Date range filtering:
GET /api/v1/posts?start_date=2023-01-01&end_date=2023-12-31&limit=20

5. Author filtering:
GET /api/v1/posts?created_by=user123&limit=10

6. Tag filtering:
GET /api/v1/posts?tags=golang,webdev,api&limit=15

7. Combined filters:
GET /api/v1/posts?search=api&sort_by=created_at&sort_order=desc&limit=20&offset=10

8. Filter by publication status:
GET /api/v1/posts?published=true&limit=10
GET /api/v1/posts?published=false&limit=10

Available sort fields:
- id
- title
- created_at (default)
- updated_at
- view_count
- like_count

Available sort orders:
- asc
- desc (default)

The enhanced query filter supports:
✓ Search in title and body fields
✓ Sort by multiple criteria
✓ Date range filtering
✓ Author/user filtering
✓ Tag filtering
✓ Publication status filtering
✓ Pagination with limit and offset
✓ Combined filtering (multiple filters together)
*/
