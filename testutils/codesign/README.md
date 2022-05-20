# codesign

We can check binary signatures only on respective platforms as of now. For example, Windows OS binary can be checked by running the `Verify` function in a Windows environment, same for MacOS. In Linux OS, we don't have any binary signature verification required

In MacOS, we use `spctl` tool
In WindowsOS, we use [`signtool`](https://docs.microsoft.com/en-us/windows/win32/seccrypto/signtool) tool

In future, we will throw errors when trying to check binary signature in a cross platform - for example Windows OS binary checking happening in MacOS
