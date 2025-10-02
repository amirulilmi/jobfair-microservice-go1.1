# User Profile Service API Collection

Base URL: http://localhost:8083
Auth: Bearer Token (dari auth-service)

## 1. Health Check

**No Auth Required**

```
GET /health
```

## 2. Profile Endpoints

### Create Profile
```
POST /api/v1/profiles
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "full_name": "John Doe",
  "phone_number": "081234567890"
}
```

### Get Profile
```
GET /api/v1/profiles
Authorization: Bearer {{token}}
```

### Get Complete Profile (with relations)
```
GET /api/v1/profiles/full
Authorization: Bearer {{token}}
```

### Update Profile
```
PUT /api/v1/profiles
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "bio": "Experienced Mining Engineer with 5+ years",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "date_of_birth": "1995-05-15T00:00:00Z",
  "gender": "male",
  "linkedin_url": "https://linkedin.com/in/johndoe"
}
```

### Get Completion Status
```
GET /api/v1/profiles/completion
Authorization: Bearer {{token}}
```

## 3. Work Experience

### Create Work Experience
```
POST /api/v1/work-experiences
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "company_name": "PT Bumi Mineral Sejahtera",
  "job_position": "Senior Mining Engineer",
  "start_date": "2020-01-15T00:00:00Z",
  "end_date": "2024-07-31T00:00:00Z",
  "is_current_job": false,
  "job_description": "Led a team of 15 engineers in planning and supervising open-pit mining operations."
}
```

### Get All Work Experiences
```
GET /api/v1/work-experiences
Authorization: Bearer {{token}}
```

### Update Work Experience
```
PUT /api/v1/work-experiences/{{id}}
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "company_name": "PT Bumi Mineral Sejahtera",
  "job_position": "Lead Mining Engineer",
  "start_date": "2020-01-15T00:00:00Z",
  "is_current_job": true,
  "job_description": "Updated description..."
}
```

### Delete Work Experience
```
DELETE /api/v1/work-experiences/{{id}}
Authorization: Bearer {{token}}
```

## 4. Education

### Create Education
```
POST /api/v1/educations
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "university": "Institut Teknologi Bandung",
  "major": "Teknik Pertambangan",
  "degree": "Bachelor",
  "start_date": "2014-08-01T00:00:00Z",
  "end_date": "2018-07-31T00:00:00Z",
  "is_current": false,
  "gpa": 3.98
}
```

### Get All Educations
```
GET /api/v1/educations
Authorization: Bearer {{token}}
```

## 5. Certifications

### Create Certification
```
POST /api/v1/certifications
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "certification_name": "Certified Mine Manager (CMM)",
  "organizer": "Ministry of Energy and Mineral Resources",
  "issue_date": "2022-03-15T00:00:00Z",
  "credential_id": "CMM-2022-12345",
  "credential_url": "https://example.com/cert/12345"
}
```

### Get All Certifications
```
GET /api/v1/certifications
Authorization: Bearer {{token}}
```

## 6. Skills

### Add Single Skill
```
POST /api/v1/skills
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "skill_name": "Mine Planning & Design",
  "skill_type": "technical",
  "proficiency_level": "expert",
  "years_of_experience": 5
}
```

### Add Multiple Skills (Bulk)
```
POST /api/v1/skills/bulk
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "technical_skills": [
    {
      "skill_name": "Mine Planning & Design",
      "skill_type": "technical",
      "proficiency_level": "expert",
      "years_of_experience": 5
    },
    {
      "skill_name": "Drilling & Blasting",
      "skill_type": "technical",
      "proficiency_level": "advanced",
      "years_of_experience": 4
    }
  ],
  "soft_skills": [
    {
      "skill_name": "Leadership",
      "skill_type": "soft",
      "proficiency_level": "advanced"
    },
    {
      "skill_name": "Problem Solving",
      "skill_type": "soft",
      "proficiency_level": "expert"
    }
  ]
}
```

### Get All Skills
```
GET /api/v1/skills
Authorization: Bearer {{token}}
```

## 7. Career Preference

### Create/Update Career Preference
```
POST /api/v1/career-preference
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "is_actively_looking": true,
  "expected_salary_min": 15000000,
  "expected_salary_max": 25000000,
  "salary_currency": "IDR",
  "is_negotiable": true,
  "preferred_work_types": ["onsite", "hybrid"],
  "preferred_locations": ["Jakarta", "Bandung"],
  "willing_to_relocate": false,
  "available_start_date": "2025-11-01T00:00:00Z"
}
```

### Get Career Preference
```
GET /api/v1/career-preference
Authorization: Bearer {{token}}
```

## 8. Position Preferences

### Create Position Preferences (Bulk)
```
POST /api/v1/position-preferences
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "positions": [
    {
      "position_name": "Senior Mining Engineer",
      "priority": 1
    },
    {
      "position_name": "Mine Operations Manager",
      "priority": 2
    },
    {
      "position_name": "Health & Safety Supervisor",
      "priority": 3
    }
  ]
}
```

### Get Position Preferences
```
GET /api/v1/position-preferences
Authorization: Bearer {{token}}
```

## 9. CV Upload

### Upload CV
```
POST /api/v1/cv
Authorization: Bearer {{token}}
Content-Type: multipart/form-data

Body (form-data):
file: [select your CV file - PDF, DOC, or DOCX]
```

### Get CV Info
```
GET /api/v1/cv
Authorization: Bearer {{token}}
```

### Delete CV
```
DELETE /api/v1/cv
Authorization: Bearer {{token}}
```

## 10. Badges

### Get User Badges
```
GET /api/v1/badges
Authorization: Bearer {{token}}
```

---

## Complete Test Flow

1. Get JWT token from auth-service
2. Create profile
3. Add work experiences
4. Add educations
5. Add certifications
6. Add skills (bulk)
7. Set career preference
8. Set position preferences
9. Upload CV
10. Get complete profile with all relations
11. Check completion status (should be close to 100%)

---

## Environment Variables for Postman

```
token: <your-jwt-token>
base_url: http://localhost:8083
```
