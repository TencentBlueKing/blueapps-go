# Contributing to blueapps-go

BlueKing Team maintains an open attitude and welcomes like-minded developers to contribute to the project. Before
getting started, please carefully read the following guidelines.

## Code License

The [MIT LICENSE](../LICENSE) is the open-source license for blueapps-go. All contributed code will be protected under
this license. Please ensure you accept this agreement before contributing code.

## Contributing Features & Enhancements

To contribute features or enhancements to the blueapps-go project, follow these steps:

- Check existing [Issues](https://github.com/TencentBlueKing/blueapps-go/issues) for any **related or similar** functionality needs. If found, discuss within that issue.
- If no relevant issue exists, create a new issue detailing your feature requirements. The BlueKing Team regularly reviews and participates in discussions.
- If the team approves the feature, supplement the issue with design details, implementation specifics, and test cases.
- Complete feature development following the [BlueKing Development Specifications](https://bk.tencent.com/docs/document/7.0/250/46218), including unit tests and documentation.
- First-time contributors must sign the [Tencent Contributor License Agreement](https://bk-cla.bktencent.com/TencentBlueKing/blueapps-go).
- Submit a [Pull Request (PR)](https://github.com/TencentBlueKing/blueapps-go/pulls) to the `main` branch, linking to the corresponding issue. PRs must include code, documentation, and unit tests.
- The BlueKing Team will promptly review the PR and merge it into `main` upon approval.

> Note: For large features, we recommend splitting requirements and submitting multiple PRs for review when possible.
> This approach maintains functionality while accelerating the review process.

## Getting Started

Developers, please refer to the document [Development Guide (In Chinese)](DEVELOP_GUIDE.md) to set up your local development environment and then start writing code.

## Git Commit Conventions

We recommend using **concise yet precise** commit messages. Follow this format:

```
git commit -m 'tag: Brief summary of changes'
```

Example:

```shell
git commit -m 'fix: Correct time display error'
```

### Tag Reference

| Tag      | Description                 |  
|----------|-----------------------------|  
| feat     | New feature/functionality   |  
| fix      | Bug fixes                   |  
| docs     | Documentation updates       |  
| style    | Formatting/comments cleanup |  
| refactor | Code restructuring          |  
| perf     | Performance optimizations   |  
| test     | Unit test modifications     |  
| chore    | Build/task adjustments      |  

## Submitting Pull Requests

If you’re working on an existing issue:

1. Comment on the issue to signal your progress and avoid duplicate efforts.
2. Fork the `main` branch to your repository.
3. Create a feature/fix branch (e.g., `fix_time_display`).
4. Implement changes while adhering to conventions, updating docs and tests.
5. Test locally and ensure all unit tests pass.
6. Submit a PR linked to the issue.

PRs should include all relevant updates: code, documentation, and usage notes.

> Note: Ensure PR titles follow commit conventions. Minimize commits within a single PR.

## Issue Reporting

We use [Issues](https://github.com/TencentBlueKing/blueapps-go/issues) to track bugs and feature requests.

When reporting bugs:

* Search for existing issues to avoid duplicates.
* For new bugs, include:
    - OS/environment details (OS, language version)
    - Version/commit ID
    - Relevant logs (exclude sensitive data)
    - Precise reproduction steps (reproducible scripts/tools preferred)

---

Best regards,  
BlueKing (PaaS) Team