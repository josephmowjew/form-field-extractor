# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2025-01-30
### Changed
- Removed redundant main.go from root directory
- All functionality now properly organized in pkg/scrapper and cmd/scrapper

## [0.1.0] - 2025-01-30
### Added
- Initial release
- PDF form field extraction using pdfcpu
- HTML form field extraction using Rod
- CLI tool for form field extraction
- Configurable timeouts and retry attempts
- Support for various input types (text, select, textarea, etc.)
- Intelligent label detection for HTML forms
- Smart field name cleaning and normalization
- JSON output format
- Comprehensive documentation

[0.1.1]: https://github.com/josephmowjew/form-field-extractor/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/josephmowjew/form-field-extractor/releases/tag/v0.1.0 