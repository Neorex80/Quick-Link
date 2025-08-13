# 🐛 GitHub Issues for QuickLink URL Shortener

This document outlines the GitHub issues to be created for feature enhancements.

## 📋 Issue #1: QR Code Generation for Short URLs

### Title
🔲 Add QR Code Generation for Short URLs

### Description
```markdown
## Feature Request: QR Code Generation

### Summary
Add QR code generation functionality to allow users to easily share shortened URLs via QR codes.

### Problem
Users often need to share URLs in physical formats or mobile-friendly ways. QR codes provide an easy way to bridge digital and physical sharing.

### Proposed Solution
- Add QR code generation for each shortened URL
- Display QR code alongside the shortened URL in the web interface
- Provide API endpoint to get QR code for any short URL
- Use a lightweight QR code library for Go

### Acceptance Criteria
- [ ] QR code is generated automatically when URL is shortened
- [ ] QR code is displayed in the web interface
- [ ] QR code can be downloaded/saved
- [ ] API endpoint `/qr/{shortCode}` returns QR code image
- [ ] QR codes are properly sized and readable
- [ ] Mobile-friendly QR code display

### Technical Implementation
- Use `github.com/skip2/go-qrcode` library
- Add QR code generation to URL shortening flow
- Update web interface to display QR codes
- Add new API endpoint for QR code retrieval

### Labels
- `enhancement`
- `feature`
- `good first issue`

### Assignees
- @Neorex80

### Milestone
- v1.1.0
```

## 📋 Issue #2: Custom Short Codes

### Title
⚙️ Add Custom Short Code Support

### Description
```markdown
## Feature Request: Custom Short Codes

### Summary
Allow users to specify custom short codes instead of only using randomly generated ones.

### Problem
Users sometimes want memorable or branded short codes (e.g., `company-name`, `event-2025`) instead of random alphanumeric codes.

### Proposed Solution
- Add optional `custom_code` parameter to the API
- Validate custom codes for uniqueness and format
- Fallback to random generation if custom code is taken
- Update web interface with custom code input option

### Acceptance Criteria
- [ ] API accepts optional `custom_code` parameter
- [ ] Custom codes are validated (alphanumeric, 3-20 characters)
- [ ] Duplicate custom codes are rejected with proper error
- [ ] Web interface has optional custom code input field
- [ ] Custom codes follow URL-safe character restrictions
- [ ] Reserved words are blocked (admin, api, www, etc.)

### Technical Implementation
- Update `ShortenRequest` struct to include `CustomCode` field
- Add validation function for custom codes
- Check uniqueness before storing
- Update web interface with optional input field
- Add proper error handling and user feedback

### Labels
- `enhancement`
- `feature`
- `api`

### Assignees
- @Neorex80

### Milestone
- v1.2.0
```

## 🔄 Implementation Workflow

### Branch Strategy
1. `feature/qr-code-generation` - QR code feature
2. `feature/custom-short-codes` - Custom codes feature
3. Merge both to `main` via Pull Requests

### Commit Message Convention
- `feat: add QR code generation for short URLs`
- `feat: implement custom short code support`
- `docs: update README with new features`
- `test: add tests for QR code generation`
- `test: add tests for custom short codes`

### Pull Request Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests pass locally
- [ ] Manual testing completed
- [ ] No breaking changes

## Screenshots (if applicable)
Add screenshots of UI changes

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added/updated
```

## 📊 GitHub Activity Plan

### Week 1: QR Code Feature
- Day 1: Create issue, create branch
- Day 2-3: Implement QR code generation
- Day 4: Update UI and documentation
- Day 5: Create PR, review, merge

### Week 2: Custom Codes Feature
- Day 1: Create issue, create branch
- Day 2-3: Implement custom code validation
- Day 4: Update UI and API
- Day 5: Create PR, review, merge

### Result
- 2 issues created and closed
- 2 feature branches
- 2 pull requests
- Multiple commits with good messages
- Updated documentation
- Enhanced application features

This will create excellent GitHub activity showing professional development workflow!
