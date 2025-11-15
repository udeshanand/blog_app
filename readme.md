
  <h1>Go Blog App</h1>
  <p>A simple and efficient blog backend API built using Go (Golang), GORM, MySQL, and HTML templates.</p>

  <h2>Features</h2>
  <ul>
    <li>User registration and authentication</li>
    <li>CRUD operations for blog posts</li>
    <li>Comment system</li>
    <li>HTML templates for rendering pages</li>
    <li>Database support with MySQL and GORM</li>
    <li>Session management with secure cookies</li>
  </ul>


  <h2>roject Structure</h2>
  <pre>
blog-app/
│
├── auth/                # Authentication logic
├── config/            # Environment configuration and DB setup
├── handlers/            # HTTP handlers
├── models/              # GORM models
├── templates/           # HTML templates rendered by Go
├── main.go              # Entry point
└── .env.example         # Sample environment file
  </pre>

  <h2>API Endpoints</h2>
  <ul>
    <li>POST <code>/register</code> – Register new user</li>
    <li>POST <code>/login</code> – User login</li>
    <li>GET <code>/posts</code> – List all posts</li>
    <li>POST <code>/posts</code> – Create new post </li>
    <li>GET <code>/posts/:id</code> – Get post by ID</li>
    <li>PUT <code>/posts/:id</code> – Update post</li>
    <li>DELETE <code>/posts/:id</code> – Delete post</li>
  </ul>
