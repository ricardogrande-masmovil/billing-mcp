# Pull Request Description Instructions

## Structure
1.  **Purpose:** What & why? Link issues (e.g., "Closes #123").
2.  **Changes Made:** Key changes (bullet points).
3.  **How to Test:** Steps for verification.
4.  **Screenshots/GIFs (if applicable).**
5.  **Checklist (Optional):** Docs updated, tests added, etc.

## Tone
- Clear, concise, professional.

## Example
**Purpose:**
This PR introduces invoice archiving. Closes #78.

**Changes Made:**
- Added `ArchiveInvoice` to `InvoiceService`.
- Implemented GORM soft-delete.
- Exposed `archiveInvoice` MCP tool.
