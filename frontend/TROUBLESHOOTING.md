# Troubleshooting Guide

## Missing Authorization Header

If you're getting "Missing authorization header" errors, follow these steps:

### 1. Check if Token is Stored

Open browser DevTools (F12) → Application/Storage → Local Storage → Check for `access_token`

### 2. Verify Token Format

The token should be a JWT string. Check in console:
```javascript
localStorage.getItem('access_token')
```

### 3. Clear and Re-login

If tokens are corrupted:
```javascript
// In browser console
localStorage.clear()
// Then login again
```

### 4. Check Network Tab

In DevTools → Network tab:
- Check if `Authorization: Bearer <token>` header is present in requests
- Verify the token value is correct

### 5. Verify Backend Services

Ensure all backend services are running:
```bash
# Check if services are up
curl http://localhost:3001/health
curl http://localhost:3002/health
# etc.
```

### 6. CORS Issues

If you see CORS errors:
- Verify backend services have CORS enabled
- Check that frontend URL is allowed in backend CORS config

### 7. Token Expiration

If token expired:
- The app should auto-refresh using refresh token
- If refresh fails, you'll be redirected to login
- Check browser console for refresh errors

## Common Issues

### Issue: "Invalid or expired token"

**Solution:**
1. Clear localStorage
2. Login again
3. Check if backend JWT secret matches

### Issue: "User not found" after login

**Solution:**
1. Verify user was created in auth service
2. Check if user exists in user service database
3. Ensure both services are connected to correct databases

### Issue: API calls return 401

**Solution:**
1. Check if token is in localStorage
2. Verify token format (should start with `eyJ`)
3. Check backend logs for token validation errors
4. Ensure JWT_SECRET matches between frontend and backend

### Issue: CORS errors

**Solution:**
Backend services already have CORS enabled. If you still see errors:
1. Check backend service logs
2. Verify request origin matches CORS config
3. Check browser console for specific CORS error message

## Debug Steps

1. **Check localStorage:**
   ```javascript
   console.log('Access Token:', localStorage.getItem('access_token'))
   console.log('Refresh Token:', localStorage.getItem('refresh_token'))
   ```

2. **Check Auth Store:**
   ```javascript
   // In React component
   const { accessToken, isAuthenticated } = useAuthStore()
   console.log('Auth State:', { accessToken, isAuthenticated })
   ```

3. **Test API directly:**
   ```bash
   # Get token from localStorage, then:
   curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:3002/api/v1/users
   ```

4. **Check Network Requests:**
   - Open DevTools → Network
   - Look for failed requests (red)
   - Check Request Headers for Authorization
   - Check Response for error details

## Still Having Issues?

1. Check browser console for errors
2. Check backend service logs
3. Verify all services are running
4. Clear browser cache and localStorage
5. Try in incognito/private window
6. Check if issue is specific to one endpoint or all endpoints

