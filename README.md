# Shared DNS Resolver for Multi-VLAN Mikrotik Setup (Proof of Concept)

> ⚠️ **Proof of Concept**  
> This project is experimental. Use at your own risk. No guarantees are provided regarding security, stability, or future compatibility.

## 🧩 Problem Statement

In a network with two separate VLANs - let's call them `main` and `secondary` - each VLAN uses a different method of accessing the internet. Despite this, the goal was to configure **a shared DNS resolver/cache** that could serve both VLANs.

The setup is built on top of the **Mikrotik hAP ax² router**.

Mikrotik's built-in DNS resolver service (as of RouterOS version `7.20beta2`) supports binding to only one interface IP. This limitation makes it unusable for a multi-VLAN environment when the resolver is configured only for the `main` VLAN.

## ❌ Solution Attempt #1 — Firewall Routing (Failed)

An initial attempt was made to route DNS traffic from the `secondary` VLAN to the resolver in the `main` VLAN using:
- Firewall rules
- Connection markers
- Multiple routing tables
- Selective inter-VLAN traffic allowances

This setup proved too complex and ultimately **failed** to provide a reliable solution.
Also, it does not give you an opportunity to cross-compile and play with containers in Mikrotik.

## ✅ Solution #2 — UDP Proxy via Container (Current Solution)

Starting in RouterOS version `7.20`, it is possible to assign **multiple veth (virtual Ethernet) interfaces** to a container. This opens up a cleaner solution:

- Deploy a simple **UDP proxy** container with one `veth` in each VLAN.
- Avoid routing complexity by directly bridging DNS requests across VLANs within the container

This approach successfully serves DNS traffic to both VLANs using a shared caching resolver.

## 🚀 Deployment

To deploy the container:

1. Review the contents of `deploy.sh`.
2. Run it to generate & upload the container .tar-file to the mikrotik
3. In Mikrotik, create the container with a command like `/container/add file=go-server.tar envlists=RESOLVER name=resolver interface=veth1,veth2 start-on-boot=yes logging=yes user=0:0`

Make sure your Mikrotik router is running RouterOS **version 7.20 or newer** and is configured to support containers.

---

## 📜 License

This project is unlicensed / released into the public domain