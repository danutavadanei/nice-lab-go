# NICE DCV Web SDK - Desktop Cloud Visualization

## NICE DCV

NICE DCV is a high-performance remote display protocol. It lets you securely deliver
remote desktops and application streaming from any cloud or data center to any device,
over varying network conditions. By using NICE DCV with Amazon EC2, you can run graphics-intensive
applications remotely on Amazon EC2 instances. You can then stream the results to
more modest client machines, which eliminates the need for expensive dedicated workstations.

## NICE DCV Web Client SDK

The NICE DCV Web Client SDK is a JavaScript library that you can use to develop your
own NICE DCV web browser client applications. Your end users can use these applications
to connect to and interact with a running NICE DCV session.

This library is exported in both ESM and UMD formats within the corresponding folders.

## NICE DCV Web UI SDK

The NICE DCV Web UI SDK is a JavaScript library exposing a single React component that
provides the user interface to interact with the NICE DCV Web Client in-session experience.

Users are still responsible for the authentication against the NICE DCV Server,
before connecting to a NICE DCV session and use the `DCVViewer` React component.
Hence, users of this library should also import the NICE DCV Web Client SDK in their React applications.