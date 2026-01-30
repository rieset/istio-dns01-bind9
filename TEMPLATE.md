# Kubernetes Operator Template Project

## Project Purpose

This project is a template for quickly starting a Kubernetes operator project with pre-configured description files and instructions for AI assistants.

## Project Goal

The project is designed to prepare a set of files and instructions for quick start of a Kubernetes operator project. It contains templates for description files, instructions, and documentation for:

- Code quality description
- Code creation rules
- AI assistant instruction templates

## Project Structure

### `operator/` Directory

The `operator/` directory is a template for operator initialization. To create an operator, run:

```bash
cd operator
operator-sdk init --domain <domain> --repo <repo>
```

After running the command, the operator project structure will be created in the `operator/` directory, and you can start development.

### `docs/` Directory

The `docs/` directory contains documentation for all functions, algorithms, and operator logic. This documentation helps:

- Understand the architecture and logic of the operator
- Document complex algorithms and solutions
- Preserve context for future development
- Ensure consistency in documentation

### Project Files

- **README.md** - description of the operator (to be created). This file contains information about the operator itself, its purpose, requirements, and deployment instructions.

- **.cursorrules** - instructions for AI assistant for the operator to be created. Contains code writing rules, style, best practices, and other development recommendations.

- **TEMPLATE.md** (this file) - information about the project as a template preparation project. Contains context about this being a template project and how to use it.

## Using the Template

1. **Operator Initialization**: Navigate to the `operator/` directory and run the Operator SDK initialization command.

2. **Documentation Setup**: Update `README.md` with information about your operator.

3. **Rules Setup**: Update `.cursorrules` with rules specific to your operator.

4. **Documentation Creation**: Start filling the `docs/` directory with documentation for functions, algorithms, and logic of your operator.

5. **Development**: Start developing the operator following the instructions in `.cursorrules` and documentation.

## Notes

- This project is a template and does not contain a ready-made operator
- All files in the project are templates and should be adapted for a specific operator
- The `operator/` directory will be populated after running the `operator-sdk init` command
- Documentation in `docs/` should be created as the operator is developed

## Context for Template Development

This file was created to preserve context during template preparation. It describes the project's purpose as a template and helps understand the structure and usage of the project when working with it.

