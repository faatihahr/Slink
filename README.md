# Slink - URL Shortener Platform

## 📝 Deskripsi Project

**Slink** adalah platform URL shortener modern yang dikembangkan dengan teknologi web terkini. Platform ini memungkinkan pengguna untuk mempersingkat URL panjang menjadi link pendek yang mudah dibagikan, dilengkapi dengan fitur analitik dan QR code generation.

## 🚀 Fitur Utama

### 🔗 **URL Shortening**
- Generate short URL otomatis dengan kode unik
- Custom alias untuk personalisasi link (maksimal 10 karakter)
- Validasi URL real-time untuk memastikan link valid
- Redirect otomatis ke URL asli

### 📊 **Link Management**
- Dashboard interaktif untuk mengelola semua link
- Halaman Links untuk melihat dan mengatur short links
- Tracking jumlah klik untuk setiap link
- Copy link dengan satu klik
- Search functionality untuk mencari link berdasarkan URL atau alias

### 📱 **QR Code Generation**
- Generate QR code otomatis untuk setiap short link
- Download QR code dalam format PNG
- Gallery QR codes yang terorganisir
- Modal viewer untuk preview QR code lebih besar

### 👤 **User Authentication**
- Sistem registrasi dan login user
- JWT-based authentication
- Profile management
- Session management yang aman

### 🎯 **Interactive Tutorial**
- Onboarding tutorial untuk user baru
- Step-by-step guidance yang interaktif
- Auto-scroll dan auto-focus ke elemen relevan
- Progress tracking dengan visual indicators
- Skip functionality untuk user yang berpengalaman

## 🛠️ Teknologi

### **Backend (Go)**
- **Framework**: Gin Web Framework
- **Database**: Supabase (PostgreSQL)
- **Authentication**: JWT tokens
- **API**: RESTful API dengan CORS support
- **QR Generation**: Library QR code generation

### **Frontend (Next.js)**
- **Framework**: Next.js 14 dengan App Router
- **Styling**: TailwindCSS
- **State Management**: React hooks (useState, useEffect)
- **Routing**: Client-side navigation
- **TypeScript**: Full TypeScript support

### **Database & Infrastructure**
- **Database**: Supabase (PostgreSQL hosting)
- **Storage**: Cloud-based database dengan RLS policies
- **Environment**: Development environment dengan hot reload
- **API Proxy**: Next.js rewrites untuk seamless API integration

## 📁 Struktur Project

```
Slink/
├── backend/                    # Go backend application
│   ├── internal/
│   │   ├── api/                # API handlers dan routes
│   │   ├── config/             # Configuration management
│   │   ├── database/           # Database connection
│   │   ├── models/             # Data models
│   │   ├── services/           # Business logic
│   │   └── utils/              # Utility functions
│   ├── main.go                 # Application entry point
│   └── .env                    # Environment variables
├── frontend/                   # Next.js frontend
│   ├── src/
│   │   ├── app/                # App Router pages
│   │   │   ├── dashboard/      # Main dashboard
│   │   │   ├── links/          # Links management
│   │   │   ├── qr-codes/       # QR codes gallery
│   │   │   └── auth/           # Authentication pages
│   │   └── components/         # Reusable components
│   ├── public/                 # Static assets
│   └── next.config.ts          # Next.js configuration
└── database/                   # Database schemas dan migrations
```

## 🎨 UI/UX Design

### **Design Principles**
- **Clean & Modern**: Interface yang minimalis dan intuitif
- **Responsive**: Optimal di desktop dan mobile
- **Interactive**: Animasi dan transisi yang smooth
- **Accessible**: Semantic HTML dan ARIA labels

### **Key Components**
- **Dashboard**: Central hub dengan quick create dan tutorial
- **Sidebar**: Navigation yang konsisten dengan active states
- **Header**: Search bar dan user profile
- **Cards**: Link cards dengan hover effects dan actions
- **Modals**: QR code preview dan detail views

## 🔐 Security Features

- **JWT Authentication**: Token-based authentication yang aman
- **Input Validation**: Server-side validation untuk semua inputs
- **CORS Protection**: Proper CORS configuration
- **SQL Injection Prevention**: Parameterized queries via Supabase
- **XSS Protection**: Sanitized user inputs

## 📈 Analytics & Tracking

- **Click Tracking**: Hit counter untuk setiap link
- **User Analytics**: Link performance per user
- **Real-time Updates**: Immediate click count updates
- **Data Visualization**: Progress bars dan statistics

## 🔄 API Endpoints

### **Authentication**
- `POST /api/register` - User registration
- `POST /api/login` - User login
- `GET /api/profile` - Get user profile

### **URL Management**
- `POST /api/shorten` - Create short URL
- `GET /api/links` - Get user links
- `GET /:shortCode` - Redirect to original URL
- `GET /api/qr/:shortCode` - Generate QR code

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Node.js 18+
- Supabase account

### Setup
1. Clone repository
2. Configure Supabase database
3. Setup backend environment variables
4. Setup frontend environment
5. Run both services

### Backend Setup
```bash
cd backend
cp .env.example .env
# Configure .env with Supabase credentials
go run main.go
```

### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

## 🎯 Use Cases

1. **Marketing Campaigns**: Short links untuk social media dan email marketing
2. **Event Management**: QR codes untuk event check-in
3. **Analytics**: Track link performance dan user engagement
4. **Brand Customization**: Custom aliases untuk brand consistency
5. **Mobile Apps**: QR codes untuk deep linking

## 🔮 Future Enhancements

- **Advanced Analytics**: Geographic tracking dan device analytics
- **Link Expiration**: Time-based link expiration
- **Password Protection**: Password-protected short links
- **Bulk Operations**: Multiple URL shortening
- **API Integration**: Public API untuk third-party integration
- **Team Collaboration**: Multi-user workspace management

## 📊 Performance Goals

- Redirect time: < 100ms
- QR code generation: < 500ms
- API response time: < 200ms
- Mobile responsiveness: 100%

---

**Slink** - Simplifying URL sharing with modern technology and intuitive design.
