# 🔧 Auto-Create Profile Fix

## 📋 Issue Fixed

**Error:**
```json
{
    "success": false,
    "message": "ERROR: insert or update on table \"career_preferences\" violates foreign key constraint \"career_preferences_profile_id_fkey\" (SQLSTATE 23503)",
    "code": "SAVE_FAILED"
}
```

**Root Cause:**
- Users trying to create preferences/work experience/education/etc before creating a profile
- Foreign key constraint requires profile to exist first

---

## ✅ Solution Implemented

Added **auto-create profile** functionality:
- All "Create" operations now automatically create profile if it doesn't exist
- Users no longer need to explicitly call `/api/v1/profiles` first
- Better user experience - just start adding data!

---

## 🔄 What Changed

### New Method Added
**File:** `internal/services/profile_service.go`

```go
func GetOrCreateProfile(userID uint) (*models.Profile, error)
```

This method:
1. Tries to get existing profile
2. If not found, automatically creates a new empty profile
3. Returns the profile (existing or new)

### Updated Services
All "Create" methods in these services now use `GetOrCreateProfile`:

✅ **PreferenceService**
- `CreateOrUpdateCareerPreference()` 
- `CreatePositionPreferences()`

✅ **WorkExperienceService**
- `Create()`

✅ **EducationService**
- `Create()`

✅ **CertificationService**
- `Create()`

✅ **SkillService**
- `Create()`
- `CreateBulk()`

✅ **CVService**
- `Upload()`

---

## 🧪 Testing

### Test 1: Create Career Preference Directly

**Before** (would fail):
```bash
# ERROR: profile not found or foreign key violation
```

**After** (works!):
```bash
curl -X POST http://localhost/api/v1/career-preference \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_type": "full_time",
    "work_location": "hybrid",
    "expected_salary_min": 15000000,
    "expected_salary_max": 25000000,
    "currency": "IDR",
    "willing_to_relocate": true,
    "available_from": "2025-11-01"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Career preference saved successfully",
  "data": {
    "id": 1,
    "profile_id": 1,
    "is_actively_looking": true,
    "expected_salary_min": 15000000,
    "expected_salary_max": 25000000,
    "salary_currency": "IDR",
    "willing_to_relocate": true,
    "available_start_date": "2025-11-01T00:00:00Z"
  }
}
```

### Test 2: Add Work Experience Directly

```bash
curl -X POST http://localhost/api/v1/work-experiences \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "Tech Corp",
    "job_position": "Software Engineer",
    "start_date": "2020-01-01",
    "end_date": "2023-12-31",
    "is_current_job": false,
    "job_description": "Developed web applications"
  }'
```

### Test 3: Upload CV Directly

```bash
curl -X POST http://localhost/api/v1/cv \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@cv.pdf"
```

### Test 4: Add Skills Directly

```bash
curl -X POST http://localhost/api/v1/skills \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "skill_name": "Python",
    "skill_type": "technical",
    "proficiency_level": "expert",
    "years_of_experience": 5
  }'
```

---

## 🚀 Deployment

### Step 1: Rebuild Service
```bash
# Stop the service
docker-compose stop user-profile-service

# Rebuild
docker-compose build user-profile-service

# Start
docker-compose up -d user-profile-service

# Check logs
docker-compose logs -f user-profile-service
```

### Step 2: Verify Fix
```bash
# Get a fresh token
TOKEN=$(curl -s -X POST http://localhost/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "SecurePass123"
  }' | jq -r '.data.access_token')

# Test creating career preference without creating profile first
curl -X POST http://localhost/api/v1/career-preference \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_type": "full_time",
    "work_location": "remote",
    "expected_salary_min": 10000000,
    "expected_salary_max": 20000000,
    "currency": "IDR",
    "willing_to_relocate": false,
    "available_from": "2025-11-01"
  }' | jq

# Should work now! ✅
```

---

## 📖 API Flow

### Old Flow (Required Manual Profile Creation)
```
1. POST /api/v1/profiles          ← Must call this first!
2. POST /api/v1/career-preference ← Then this
3. POST /api/v1/work-experiences  ← Then this
```

### New Flow (Auto-Create Profile)
```
1. POST /api/v1/career-preference ← Just start here! ✅
   (Profile auto-created if not exists)
   
2. POST /api/v1/work-experiences  ← Or here! ✅
   (Profile auto-created if not exists)

3. POST /api/v1/cv                ← Or even here! ✅
   (Profile auto-created if not exists)
```

---

## 💡 Benefits

✅ **Better UX**: Users can start adding data immediately
✅ **Fewer Errors**: No more foreign key constraint violations
✅ **Flexible**: Can create profile data in any order
✅ **Backwards Compatible**: Explicit profile creation still works

---

## 🔍 Technical Details

### Profile Creation Logic

```go
func (s *profileService) GetOrCreateProfile(userID uint) (*models.Profile, error) {
    // Try to get existing profile
    profile, err := s.profileRepo.GetByUserID(userID)
    if err == nil {
        return profile, nil  // Profile exists, return it
    }

    // If not found, create new profile
    if errors.Is(err, gorm.ErrRecordNotFound) {
        newProfile := &models.Profile{
            UserID:           userID,
            FullName:         "",
            PhoneNumber:      "",
            CompletionStatus: 0,
        }

        err = s.profileRepo.Create(newProfile)
        if err != nil {
            return nil, err
        }

        return newProfile, nil
    }

    return nil, err
}
```

### When Profile Is Created

Profile is automatically created on **first data entry**:
- First career preference
- First work experience
- First education entry
- First certification
- First skill
- First CV upload

### Profile Completion Status

- Initially: 0%
- Updates automatically as user adds more data
- Check completion: `GET /api/v1/profiles/completion`

---

## 🧪 Complete Test Script

```bash
#!/bin/bash

echo "Testing Auto-Create Profile Feature"
echo "===================================="

# Login
TOKEN=$(curl -s -X POST http://localhost/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "password123"
  }' | jq -r '.data.access_token')

echo "Token obtained: ${TOKEN:0:20}..."
echo ""

# Test 1: Create career preference (will auto-create profile)
echo "Test 1: Creating career preference..."
curl -s -X POST http://localhost/api/v1/career-preference \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_type": "full_time",
    "work_location": "hybrid",
    "expected_salary_min": 15000000,
    "expected_salary_max": 25000000,
    "currency": "IDR",
    "willing_to_relocate": true,
    "available_from": "2025-11-01"
  }' | jq

echo ""

# Test 2: Check if profile was auto-created
echo "Test 2: Checking profile..."
curl -s -X GET http://localhost/api/v1/profiles \
  -H "Authorization: Bearer $TOKEN" | jq

echo ""

# Test 3: Add work experience
echo "Test 3: Adding work experience..."
curl -s -X POST http://localhost/api/v1/work-experiences \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "Tech Corp",
    "job_position": "Software Engineer",
    "start_date": "2020-01-01",
    "is_current_job": false,
    "end_date": "2023-12-31",
    "job_description": "Developed applications"
  }' | jq

echo ""

# Test 4: Check completion status
echo "Test 4: Checking completion status..."
curl -s -X GET http://localhost/api/v1/profiles/completion \
  -H "Authorization: Bearer $TOKEN" | jq

echo ""
echo "All tests completed!"
```

---

## 🐛 Troubleshooting

### Issue: Still getting foreign key error

**Check:**
1. Service rebuilt? `docker-compose build user-profile-service`
2. Service restarted? `docker-compose restart user-profile-service`
3. Using latest code? `git pull`
4. Database migrations run? Check migrations

**Solution:**
```bash
# Full rebuild
docker-compose down
docker-compose build --no-cache user-profile-service
docker-compose up -d
```

### Issue: Profile created but empty

**Expected Behavior:**
- Profile starts with 0% completion
- Add data to increase completion %
- Update profile info later via `PUT /api/v1/profiles`

---

## 📊 Impact

**Files Modified:** 7
- `profile_service.go` (added GetOrCreateProfile)
- `preference_service.go` (updated Create methods)
- `work_experience_service.go` (updated Create)
- `education_service.go` (updated Create)
- `certification_service.go` (updated Create)
- `skill_service.go` (updated Create & CreateBulk)
- `cv_service.go` (updated Upload)

**Breaking Changes:** None ✅
- All existing functionality works as before
- Additional auto-create feature added

---

## ✨ Summary

**Before:**
```
User → Create Preference → ❌ Foreign Key Error
```

**After:**
```
User → Create Preference → ✅ Profile Auto-Created → Success!
```

No more foreign key constraint violations! 🎉

---

**Last Updated**: October 3, 2025
