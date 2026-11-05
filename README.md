# GitOps Config Manager

[![CI](https://github.com/mtcass/gitops-config-manager/actions/workflows/ci.yaml/badge.svg)](https://github.com/mtcass/gitops-config-manager/actions/workflows/ci.yaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.27%2B-blue.svg)](https://golang.org/)

A lightweight, self-service **Internal Developer Platform (IDP)** built with **Go and HTMX** to visually manage JSON configuration files and automate Git Merge Requests with zero database requirements.

---

## 💡 Why GitOps Config Manager?

In modern cloud-native environments, developers often need to update configuration JSONs, but editing files directly in Git repositories can lead to syntax errors, schema violations, and review bottlenecks. 

**GitOps Config Manager** solves this by providing:
* **Zero Database Overhead:** The Git repository remains the single source of truth.
* **Instant UI with HTMX:** A fast, reactive interface without the heavy footprint of a Single Page Application (React/Vue).
* **Automated GitOps Workflows:** Safely creates branches, commits configuration changes, and opens Merge/Pull Requests automatically via API.
* **Single Binary Deployment:** Compiled into a tiny, self-contained container ready for production.

---

## 🏗️ Architecture Overview

```text
+-------------------+           HTTP / HTMX           +--------------------+
|                   |  ---------------------------->  |                    |
|  Browser / Client |                                 |     Go Backend     |
|                   |  <----------------------------  |  (Embedded Assets) |
+-------------------+      HTML Fragments / UI        +---------+----------+
                                                                |
                                                                | SDK (GitLab / GitHub API)
                                                                v
                                                      +--------------------+
                                                      |                    |
                                                      |   Git Repository   |
                                                      |  (Configs & MRs)   |
                                                      +--------------------+
```

---

## 🧪 Companion Sandbox Repository

To test the application end-to-end without touching real production repositories, a companion sandbox repository is available:

👉 **[gitops-config-manager-sandbox](https://github.com/mtcass/gitops-config-manager-sandbox)**

### Why Use the Sandbox?
* **Zero Setup Overheads:** Contains pre-configured JSON configuration files and schemas ready for testing.
* **Safe Playground:** Allows you to test branch creation, file edits, and Pull Request generation in an isolated environment.
* **Reference Spec:** Serves as a live reference for how target configuration repositories should be structured.

## ⚙️ Getting Started & Prerequisites

- Go `1.27+` installed locally.
- A clone or fork of the **[gitops-config-manager-sandbox](https://github.com/mtcass/gitops-config-manager-sandbox)** repository.
- A Personal Access Token (GitHub or GitLab).

### GitHub Token Permissions (Fine-grained)
If using a **Fine-grained Personal Access Token** targeting your sandbox repository, grant the following **Repository permissions**:

| Permission | Access Level | Purpose |
| :--- | :--- | :--- |
| **Contents** | `Read and write` | Reads JSON files and pushes new branches with updated configs. |
| **Pull requests** | `Read and write` | Creates, tracks, and manages Pull Requests automatically. |
| **Metadata** | `Read-only` | Required by GitHub API to fetch repository details. |
