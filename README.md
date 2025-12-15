# package-suggestions
A repository of package suggestions based on a fuzzy name.

## How to add a matcher

1. Create a new JSON file in the `matchers/` directory (e.g., `matchers/my_package.json`).
2. Add the following structure to the file:

```json
{
  "matcher": "(?i)^(?:.*)\bmy_package\b",
  "test_ok": ["my_package", "I need my_package"],
  "test_fail": ["other_package", "not_my_package"],
  "packages": [{ "name": "my-system-package", "manager": "apt" }]
}
```

### Fields

- `matcher`: A Go-compatible regular expression to match the package name.
- `test_ok`: A list of strings that *must* match the regex.
- `test_fail`: A list of strings that *must not* match the regex.
- `packages`: A list of system packages to suggest if the matcher succeeds. Each package object contains:
    - `name`: The system package name.
    - `manager`: The package manager (e.g., `apt`, `dnf`, `brew`).

### Testing and Regeneration

To test your new matcher and regenerate the registry, run:

```bash
go run main.go > registry.json
```

This command validates all matchers against their test cases and outputs the
combined registry to `registry.json`. If any test case fails, the command will
exit with an error log.
