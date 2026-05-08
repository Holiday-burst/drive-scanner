This is a community tool, not an official Ubiquiti utility. I am sharing
  both the source code and a pre-compiled binary so you can verify what
  it does before running it. You have two options:

  ────────────────────────────────────────────────────────────────────────
 
  Option A — Compile from source yourself (most transparent)
  ────────────────────────────────────────────────────────────────────────

  The full source (124 lines of Go) is here:
    [https://github.com/<your-handle>/drive-scanner/blob/main/driveScanner.go](https://github.com/Holiday-burst/drive-scanner/blob/main/driveScanner.go)

  Read it first — it's short and there's nothing hidden.

Getting Started
To begin, open your terminal and clone the repository using the command git clone https://github.com/Holiday-burst/drive-scanner. Once the download is complete, navigate into the project folder by typing cd drive-scanner.

Building the Tool
Before compiling, you need to prepare the Go environment. If the project doesn't already have a configuration file, initialize it by running go mod init drive-scanner. After initialization, you can build the executable file by executing go build -o driveScanner driveScanner.go. This will generate a binary file named driveScanner in your current directory.

Deploying to NAS
To move the tool to your NAS, I recommend using the scp command. This allows you to securely transfer the file over your network. The general syntax for this is scp ./driveScanner [username]@[nas-ip-address]:/[target-directory].

For example, if you want to send the file to the temporary folder of a NAS with the IP address 192.168.1.100 using the root account, you would run:

scp ./driveScanner root@192.168.1.100:/tmp

Final Execution
After the transfer is complete, log in to your NAS via SSH. Make sure to grant the file execution permissions by running chmod +x /tmp/driveScanner. You can then start scanning your drives by executing the program followed by the directory path you wish to analyze.

  ────────────────────────────────────────────────────────────────────────
 
  Option B — Use my pre-built binary
  ────────────────────────────────────────────────────────────────────────

  Download:
    wget -O /tmp/driveScanner https://github.com/Holiday-burst/drive-scanner/releases/download/a001/driveScanner_arm64

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
    - Prints any file that triggers a panic
    - With -quarantine, moves problem files out (no deletion)

  It does not:
    - Make network calls
    - Modify system configuration
    - Delete any files
    - Read anything outside the directories you give it
    - Send telemetry

  Hope this helps. Let me know if you have questions about the source.







  wget -O /tmp/driveScanner https://github.com/Holiday-burst/drive-scanner/releases/download/a003/scan-bplist-dos
