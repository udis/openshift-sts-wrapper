# Clarifying Questions for OpenShift STS Installation Wrapper

Please answer the following questions to help me create a comprehensive PRD. You can edit this file directly with your answers.

## 1. Target Users & Use Cases

**Q1.1:** Who is the primary user of this wrapper tool?
- A) DevOps engineers with moderate OpenShift experience
- B) Junior developers new to OpenShift
- C) Platform engineers managing multiple clusters
- D) All of the above

**Your answer:** D

**Q1.2:** What is the main pain point this wrapper solves?
- A) Too many manual steps prone to errors
- B) Complexity of extracting binaries and credentials
- C) Difficulty remembering the correct sequence
- D) All of the above

**Your answer:** D

## 2. Scope & Functionality

**Q2.1:** Should this wrapper be:
- A) A single script that runs all steps sequentially
- B) A CLI tool with subcommands for each major step (e.g., `wrapper extract`, `wrapper create-manifests`, etc.)
- C) An interactive tool that prompts the user for inputs step-by-step
- D) A combination of B and C (CLI with optional interactive mode)

**Your answer:** D. I'd like a tool that is _also_ able to restart from an already started process, stopped mid way, but if the complexity is too grow, it's ok to have C

**Q2.2:** Which steps from the STS installation process should be included? (Select all that apply)
- A) Extract credentials requests from release image
- B) Extract openshift-install and ccoctl binaries
- C) Create install-config.yaml (interactive or template-based)
- D) Create manifests
- E) Run ccoctl to create AWS resources
- F) Copy generated manifests and TLS files to correct locations
- G) Run openshift-install to create cluster
- H) Post-install verification
- I) All of the above

**Your answer:** I

**Q2.3:** For the install-config.yaml generation, should the tool:
- A) Use an existing template file that users modify
- B) Interactively prompt for all required fields
- C) Accept configuration via command-line arguments/flags
- D) Support multiple methods (e.g., B and C)

**Your answer:** just use the command defined in sts.md (`openshift-install create install-config`)

## 3. Pull Secret Automation

**Q3.1:** For the pull-secret automation, what level of automation is acceptable?
- A) Fully automated (use browser automation/headless browser to login and download)
- B) Semi-automated (open browser to the page, user authenticates manually, tool detects and extracts the secret)
- C) Guided manual (tool provides instructions and validates the file format)
- D) Optional automation (default to manual, but support automation if credentials are provided)

**Your answer:** Just make an online search to ensure some possibility to automation exists. If not (and be honest) just prompt the user that "we need to get the pull-secret from the link" and wait for user permission to continue after they take the file manually

**Q3.2:** If implementing browser automation, should it:
- A) Use existing authentication tokens/cookies if available
- B) Prompt for Red Hat account credentials (username/password)
- C) Support SSO/OAuth flow
- D) Not implement this feature initially (mark as future enhancement)

**Your answer:** just open the browser from CLI (xdg-open maybe)

## 4. Error Handling & Validation

**Q4.1:** How should the tool handle errors during execution?
- A) Stop immediately on any error
- B) Continue where possible, report errors at the end
- C) Prompt user whether to continue after each error
- D) Support both A and B via a flag (--strict mode)

**Your answer:** C

**Q4.2:** Should the tool validate prerequisites before starting? (Select all that apply)
- A) Check AWS credentials are configured
- B) Verify required CLI tools are installed (oc, aws cli)
- C) Validate the release image is accessible
- D) Check AWS permissions are sufficient
- E) All of the above

**Your answer:** B, C

## 5. Output & Logging

**Q5.1:** What level of output/logging should the tool provide?
- A) Minimal (only errors and critical steps)
- B) Normal (progress indicators for each major step)
- C) Verbose (detailed output from all commands)
- D) Configurable via flags (--quiet, --verbose, etc.)

**Your answer:** D

**Q5.2:** Should the tool:
- A) Save logs to a file automatically
- B) Display logs to console only
- C) Both (console + optional log file)

**Your answer:** B

## 6. Configuration & Flexibility

**Q6.1:** Should the tool support both public S3 bucket and private S3 bucket (CloudFront) options?
- A) Yes, via a flag (e.g., `--private-bucket`)
- B) No, default to public bucket only
- C) Yes, but ask interactively during execution

**Your answer:** A and C

**Q6.2:** How should users specify required inputs (cluster name, region, etc.)?
- A) Command-line flags only
- B) Configuration file (YAML/JSON)
- C) Environment variables
- D) Multiple methods supported (priority: flags > config file > env vars > interactive prompts)

**Your answer:** D

## 7. Implementation Language & Distribution

**Q7.1:** What language/technology should this wrapper be implemented in?
- A) Bash script (simple, no dependencies)
- B) Python (better error handling, easier to maintain)
- C) Go (single binary, no runtime dependencies)
- D) Your preference

**Your answer:** C

**Q7.2:** How should the tool be distributed?
- A) Single executable/script in a Git repository
- B) Package for package managers (pip, npm, etc.)
- C) Container image
- D) Multiple options

**Your answer:** N/A

## 8. Success Criteria

**Q8.1:** How will we measure the success of this tool?
- A) Reduces installation time from ~X minutes to ~Y minutes
- B) Reduces installation errors by X%
- C) Adoption by X users/teams within Y timeframe
- D) All of the above (please specify metrics if possible)

**Your answer:** just make installation easier

## 9. Non-Goals (Out of Scope)

**Q9.1:** What should this tool explicitly NOT do? (Select all that apply)
- A) Manage/destroy existing clusters
- B) Handle cluster upgrades
- C) Configure post-installation settings (beyond verification)
- D) Support non-AWS cloud providers
- E) Other (please specify below)

**Your answer:** All the above

## 10. Additional Context

**Q10.1:** Are there any existing tools or scripts this should integrate with or replace?

**Your answer:** NO

**Q10.2:** Any other requirements, constraints, or considerations I should know about?

**Your answer:** None that I can think of now

---

**Instructions:** Please fill in your answers above and let me know when you're ready. I'll use your responses to generate the PRD.
