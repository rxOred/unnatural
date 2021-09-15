# Unnatural

Unnatural is a an elf anomaly detector written in golang for linux systems.
It is capable of detecting malware infected binaries and disinfect those binaries to some extent

Unnatural also comes with a runtime mode, which is capable of scanning a process for
suspesicous content at runtime. it uses libbpfgo tracing library to achieve this.
