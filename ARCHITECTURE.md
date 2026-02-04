
ARCHITECTURE.md
Huggkey: The Sovereign Identity & Transport Mesh
> Status: Experimental / Pre-Alpha
> Author: Nathaniel Anthony (nathfavour)
> Core Stack: Go, libp2p, Ed25519, FUSE
> 
1. The Manifesto (Why this exists)
The current internet identity stack is broken for builders.
 * Identity is Rental: Passkeys and OAuth tokens are owned by platforms (Apple, Google), not users.
 * Connectivity is Centralized: Connecting two personal devices usually requires a third-party server or complex VPNs.
 * Agents are Second-Class: AI agents (like Auracrab) lack a verifiable, scoped identity to perform autonomous actions securely.
Huggkey is a Sovereign Identity & Transport Mesh. It combines cryptographic ownership (Identity) with peer-to-peer networking (Transport) to create a "Trust Island" where your devices—laptops, VPSs, phones, and AI agents—recognize and talk to each other directly.
2. System Overview
Huggkey operates on a "Layered Sovereignty" model. It does not rely on a central backend.
The Stack
| Layer | Component | Name | Technology |
|---|---|---|---|
| L4: Application | FUSE, SSH, Agents | The Service | go-fuse, Custom Proxy |
| L3: Transport | NAT Traversal, Streams | Einstein Hole | go-libp2p, QUIC, Noise |
| L2: Trust | Handshake, Delegation | The Hugg | Ed25519 Signatures, CapTokens |
| L1: Identity | Keys, Storage, Vault | HuggID | Argon2id, AES-GCM, Seed Phrases |
3. Core Components
3.1 Layer 1: Identity (HuggID)
Unlike Passkeys, a HuggID is a portable mathematical proof.
 * Format: hugg:v1:<base58_public_key>
 * Storage: A local JSON vault (~/.config/hugg/identity.vault) encrypted at rest.
 * Root Secret: A 24-word BIP-39 mnemonic phrase. This allows for paper backups and cross-device restoration without a cloud sync.
 * Protection: The private key is unlocked via:
   * Headless (Server/VPS): Environment variable or TPM.
   * Desktop (Arch/KDE): System Keyring / Biometrics (via fprintd).
3.2 Layer 2: Trust (The Hugg Protocol)
This layer replaces the "Login" screen.
The Handshake Protocol:
 * Device A (Client) connects to Device B (Host) via the Einstein Hole.
 * Device B sends a cryptographic Challenge (32-byte nonce).
 * Device A signs the nonce with its HuggID private key.
 * Device B verifies the signature against its "Trusted Peers" list (known_huggs).
 * Result: If valid, the stream is upgraded to an authenticated session.
3.3 Layer 3: Transport (Einstein Hole)
This is the networking engine. It ensures connectivity regardless of network topology.
 * Discovery: Uses a private DHT (Distributed Hash Table) topic to find peers based on their HuggID.
 * NAT Traversal:
   * Strategy A: UDP Hole Punching (Stun).
   * Strategy B: UPnP (Local Network).
   * Strategy C: Relay (Fallback to a user-owned VPS relay daemon).
 * Encryption: All traffic is wrapped in the Noise Protocol Framework (similar to WireGuard) to prevent metadata leakage to ISPs.
3.4 Layer 4: Application (The Interfaces)
Huggkey exposes its power via specific interfaces:
 * Wormhole-FS: A FUSE mount that presents remote directories as local folders (~/hugg/vps-01).
 * HuggSSH: A wrapper that injects the HuggID signature into the SSH handshake.
 * Agent Proxy: A local socket for AI agents to request scoped signatures.
4. Key Innovation: Agent Sovereignty
This architecture is explicitly designed for the Agent Economy.
Problem
An AI agent (e.g., Auracrab) running on a background daemon needs to deploy code to a VPS. It cannot "type" a password or use FaceID.
Solution: Cryptographic Delegation (Macaroons)
Huggkey implements Scoped Sub-Identities.
 * Issue: The Master HuggID signs a capability token (Token A).
   * Scope: allow:filesystem:read
   * Target: /var/www/project-tenchat
   * Expiry: 3600s
 * Delegate: This token is given to the Auracrab daemon.
 * Act: Auracrab presents Token A to the VPS.
 * Verify: The VPS verifies the Master's signature on the token. It grants access only to the specified scope.
5. Data Flow: The "Einstein Link"
When a user runs hugg link <vps-peer-id>:
 * Lookup: The CLI queries the DHT for <vps-peer-id>.
 * Pathfinding: libp2p determines the best route (Local LAN -> QUIC Hole Punch -> Relay).
 * Tunneling: A secure, multiplexed stream (Yamux) is established.
 * Auth: The Hugg Protocol runs the handshake.
 * Mount: If successful, the remote filesystem is mounted locally via FUSE.
6. Security Model
 * Zero-Trust: Every request is authenticated. There is no "internal network."
 * Perfect Forward Secrecy: Session keys are ephemeral and rotated per connection.
 * Attack Surface:
   * Compromised Laptop: Attacker gets the encrypted vault. Without the passphrase/biometric, it is useless.
   * Compromised VPS: Attacker gets the public keys of peers. They cannot impersonate the user.
 * Recovery:
   * Seed Phrase: Restore identity on a new machine.
   * Social Recovery (Planned): M-of-N secret sharing with trusted HuggIDs.
7. Integration Roadmap
 * [ ] Phase 1: The Core. Build huggd (daemon) and hugg (CLI) with basic Identity + P2P Ping.
 * [ ] Phase 2: The Hole. Implement robust NAT traversal and the Einstein Hole transport layer.
 * [ ] Phase 3: The Vault. Integrate with Whisperrkeep for P2P secret syncing.
 * [ ] Phase 4: The Agent. Build the delegation API for Auracrab.
8. Directory Structure
/cmd
  /hugg        # User CLI
  /huggd       # System Daemon
/pkg
  /identity    # Ed25519 & Argon2id logic
  /hole        # Libp2p host & NAT traversal
  /protocol    # Handshake & Wire format
  /agent       # Macaroon/Token delegation
  /fuse        # Filesystem implementation
/docs
  ARCHITECTURE.md

> "We are building the nervous system for a sovereign digital body." — nathfavour
> 














### adjusted version

ARCHITECTURE.md: Huggkey
1. Executive Summary
Huggkey is a next-generation, post-quantum resilient authentication system designed to replace traditional passkeys. It leverages a Hybrid Signature Scheme to provide "Double-Action" security: combining the industry-proven speed of Ed25519 with the quantum-resistance of ML-DSA (Dilithium).
2. Design Philosophy
 * Minimalism: No "legacy bloat." We prioritize modern primitives over backwards compatibility with RSA or ECDSA.
 * Quantum-First: Built on the assumption that "Harvest Now, Decrypt Later" is a real-world threat.
 * Performance: Designed to work within the latency constraints of Bluetooth Low Energy (BLE) and high-concurrency Go environments.
3. System Architecture
A. The Cryptographic Core (internal/crypto)
We utilize a composite hybrid approach where the security of the credential is the sum of two independent mathematical problems.
 * Classical Layer: Ed25519 (Elliptic Curve Discrete Logarithm Problem).
 * Post-Quantum Layer: ML-DSA / Dilithium (Module Learning with Errors).
 * Implementation: Powered by the Cloudflare CIRCL library in Go.
B. Credential Management (pkg/auth)
Unlike standard WebAuthn which often relies on hardware enclaves, Huggkey is designed as a Software-Defined Authenticator.
 * Key Generation: Creates a bound pair: PK_{hybrid} = (PK_{ed25519}, PK_{dilithium}).
 * Signature Format: A concatenated byte-stream (approx. 2.5 KB).
 * Storage: Keys are stored using the OS-level secure storage (KDE Wallet/Arch Linux secret-service or Android Keystore).
C. The Transport Layer (BLE / Network)
To handle the "PQC Bloat" (35x larger signatures), the architecture employs:
 * Fragmentation Engine: Logic to split the 2.5 KB signature into MTU-sized packets for Bluetooth delivery.
 * Encoding: Using Protobuf or MsgPack for serialization to minimize overhead compared to JSON/Base64.
4. Sequence Diagram: Authentication Ceremony
 * Challenge: Server sends a 32-byte random salt.
 * Hybrid Sign: Device signs the salt with both private keys.
 * Proof Generation (Optional/Beta): For bandwidth-critical paths, a ZK-STARK wrapper compresses the hybrid signature.
 * Verification: Server validates both layers. If either fails, the transaction is rejected.
5. Security Analysis
| Threat | Mitigation |
|---|---|
| Classical Computing Attack | Protected by Ed25519 (128-bit security). |
| Quantum Computing Attack | Protected by ML-DSA (Lattice-based security). |
| Replay Attack | Mitigated by time-bound, random server challenges. |
| Side-Channel Attack | Dilithium's constant-time implementation prevents timing leaks. |
6. Project Structure
huggkey/
├── cmd/                # Entry points (CLI, Agent)
├── internal/
│   ├── crypto/         # Hybrid EdDilithium wrappers
│   ├── storage/        # Secure key storage drivers
│   └── transport/      # BLE/TCP fragmentation logic
├── pkg/
│   ├── api/            # Go-based Auth API
│   └── types/          # Shared protobuf/struct definitions
└── web/                # TS/WASM verification library

7. Future Roadmap
 * STARK Compression: Integrating gnark to reduce signature transmission size over low-bandwidth BLE.
 * Clarigggz Integration: Porting the verification logic to the smart glasses simulation.
 * Agentic Visas: Extending the payload to include scoped permissions for autonomous AI agents.
