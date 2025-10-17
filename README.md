# OpenShift STS Installation Wrapper

A CLI tool that automates the installation of OpenShift clusters with AWS Security Token Service (STS) authentication.

## Features

- **Automated Workflow**: Consolidates 10+ manual steps into a single command
- **Error Recovery**: Resume installations from where they left off
- **Flexible Configuration**: Support for CLI flags, config files, and environment variables
- **Interactive Guidance**: Clear progress indicators and error handling
- **Version-Aware**: Automatically handles multiple OpenShift versions

## Installation

### From Source

```bash
git clone https://github.com/clobrano/ccoctl-sso.git
cd ccoctl-sso
make build
sudo make install
```

## Prerequisites

- `oc` (OpenShift CLI) must be installed and in your PATH
- AWS credentials configured
- Pull secret from Red Hat (will be prompted if not provided)

## Usage

### Full Installation

```bash
openshift-sts-installer install \
  --release-image=quay.io/openshift-release-dev/ocp-release:4.12.0-x86_64 \
  --cluster-name=my-cluster \
  --region=us-east-2 \
  --pull-secret=./pull-secret.json
```

### With Private S3 Bucket

```bash
openshift-sts-installer install \
  --release-image=quay.io/openshift-release-dev/ocp-release:4.12.0-x86_64 \
  --cluster-name=my-cluster \
  --region=us-east-2 \
  --pull-secret=./pull-secret.json \
  --private-bucket
```

### Using a Configuration File

Create `openshift-sts-installer.yaml`:

```yaml
releaseImage: quay.io/openshift-release-dev/ocp-release:4.12.0-x86_64
clusterName: my-cluster
awsRegion: us-east-2
pullSecretPath: ./pull-secret.json
privateBucket: false
outputDir: _output
```

Then run:

```bash
openshift-sts-installer install
```

### Resume from Specific Step

If installation was interrupted:

```bash
openshift-sts-installer install --start-from-step=6
```

Step numbers:
1-2. Extract credentials requests
3. Extract binaries
4. Create install-config.yaml
5. Set credentialsMode
6. Create manifests
7. Create AWS resources
8-9. Copy files
10. Deploy cluster

### Cleanup After Failed Installation

```bash
openshift-sts-installer cleanup \
  --cluster-name=my-cluster \
  --region=us-east-2
```

## Environment Variables

You can also configure via environment variables:

```bash
export OPENSHIFT_STS_RELEASE_IMAGE=quay.io/openshift-release-dev/ocp-release:4.12.0-x86_64
export OPENSHIFT_STS_CLUSTER_NAME=my-cluster
export OPENSHIFT_STS_AWS_REGION=us-east-2
export OPENSHIFT_STS_PULL_SECRET_PATH=./pull-secret.json
export OPENSHIFT_STS_PRIVATE_BUCKET=true

openshift-sts-installer install
```

## Configuration Priority

Configuration sources are merged with the following priority (highest to lowest):

1. CLI flags
2. Configuration file
3. Environment variables
4. Interactive prompts

## Directory Structure

The tool creates the following directory structure:

```
./
├── artifacts/
│   └── 4.12.0-x86_64/       # Version-specific artifacts
│       ├── bin/              # Extracted binaries
│       ├── credreqs/         # Credentials requests
│       └── install-config.yaml
├── _output/                  # ccoctl generated files
│   ├── manifests/
│   └── tls/
├── manifests/                # Installation manifests
├── tls/                      # TLS certificates
└── pull-secret.json          # Pull secret
```

## Verbosity Control

```bash
# Quiet mode (errors only)
openshift-sts-installer install --quiet

# Verbose mode (detailed output)
openshift-sts-installer install --verbose
```

## Development

### Running Tests

```bash
make test
```

### Test Coverage

```bash
make test-coverage
```

### Code Quality

```bash
make check  # Runs fmt, vet, and test
```

### Building

```bash
make build
```

## Troubleshooting

### Pull Secret Issues

If you don't have a pull secret, the tool will:
1. Display a message with the download URL
2. Attempt to open your browser to the Red Hat portal
3. Wait for you to provide the path to the downloaded file

### Step Detection

The tool automatically detects completed steps by checking for:
- Existence of directories and files
- Content of configuration files
- Presence of artifacts

If detection fails, use `--start-from-step` to manually specify where to resume.

### AWS Permissions

The tool does not validate AWS permissions before starting. If you encounter AWS errors during execution, verify that your AWS credentials have the required permissions for:
- S3 bucket creation
- IAM role/policy creation
- OIDC provider creation

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
