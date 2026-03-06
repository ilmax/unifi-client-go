# Lessons

- Failure mode: Broke GitHub Actions YAML with a heredoc that was not properly indented, making the workflow invalid. Detection signal: GitHub reported "Invalid workflow file" at a specific line. Prevention rule: Avoid multi-line heredocs in workflow YAML; prefer `python -c` or ensure heredoc content is safely indented (or use `<<-'EOF'` with tab stripping).
