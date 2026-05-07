This is a community tool, not an official Ubiquiti utility. I am sharing
  both the source code and a pre-compiled binary so you can verify what
  it does before running it. You have two options:

  ────────────────────────────────────────────────────────────────────────
 
  Option A — Compile from source yourself (most transparent)
  ────────────────────────────────────────────────────────────────────────

  The full source (124 lines of Go) is here:
    [https://github.com/<your-handle>/drive-scanner/blob/main/driveScanner.go](https://github.com/Holiday-burst/drive-scanner/blob/main/driveScanner.go)

  Read it first — it's short and there's nothing hidden.
  Deployment and Execution Guide
  1. Clone the Repository
  Open your terminal and run the following commands to download the source code:
    git clone https://github.com/your-username/drive-scanner
    cd drive-scanner
  2. Initialize Go Module
  If the project does not already contain a go.mod file, initialize it with this command:
    go mod init drive-scanner
  3. Build the Executable
  Compile the Go source code into a binary file:
    go build -o driveScanner driveScanner.go
  4. Deploy to NAS via SCP
  To transfer the compiled binary to your NAS, use the scp command. Please refer to the following usage example:
  Usage Syntax:
    scp ./driveScanner [username]@[nas-ip-address]:/[target-directory]
  Example:
    # Example: Sending the file to the "admin" user's home folder on the NAS
    scp ./driveScanner root@192.168.1.100:/tmp

  ────────────────────────────────────────────────────────────────────────
 
  Option B — Use my pre-built binary
  ────────────────────────────────────────────────────────────────────────

  Download:
    wget -O /tmp/driveScanner \

  [https://github.com/<your-handle>/drive-scanner/releases/download/v1.0/driveScanner_arm64](https://github.com/Holiday-burst/drive-scanner/blob/main/driveScanner_arm64)

  Verify the SHA256 matches:
    echo "7bdb4c60ea819a915a0b452147dd17513f02fbed4e32ddeecc1367e37ee772db
  /tmp/driveScanner" | sha256sum -c

  If it doesn't match, do not run it.

  ────────────────────────────────────────────────────────────────────────
 
  Running the tool (Option A or B)
  ────────────────────────────────────────────────────────────────────────

    chmod +x /tmp/driveScanner

    # Step 1: scan only — read-only, just lists problem files
    sudo /tmp/driveScanner /volume/*/.srv/.unifi-drive/homes/

    # Step 2: review the output, then quarantine if you agree
    sudo /tmp/driveScanner -quarantine /tmp/drive_quarantine \
         /volume/*/.srv/.unifi-drive/homes/

    # Step 3: restart the service
    sudo systemctl restart unifi-drive

  Files moved to /tmp/drive_quarantine/ are not deleted — they are renamed
  out of the photo backup folder so unifi-drive will skip them. You can
  move them back any time.

  ────────────────────────────────────────────────────────────────────────
  
  What the tool does (verifiable from the source)
  ────────────────────────────────────────────────────────────────────────

    - Walks the photo backup folder you give it
    - For each .heic/.heif/.jpg/.jpeg/.png/.tif file, runs imagemeta.Decode
      (the same library unifi-drive uses internally) inside a recover()
    - Prints any file that triggers a panic
    - With -quarantine, moves problem files out (no deletion)

  It does not:
    - Make network calls
    - Modify system configuration
    - Delete any files
    - Read anything outside the directories you give it
    - Send telemetry

  Hope this helps. Let me know if you have questions about the source.
