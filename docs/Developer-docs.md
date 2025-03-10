# Developer docs

## Code structure

```sh
cmd/ # Executable commands
  server/ # Main web server - entry point, but no logic
internal/ # The real implementation
  features/ # Groups of features
    auth/ # Authentication logic
  server/ # Http handlers. Package contains the root HTTP handler
    views/ # HTML Views
    ioc/ # Bootstrap an initialized object
  testing/ # Test helpers
```
