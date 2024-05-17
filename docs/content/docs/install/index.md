---
title: 'Install'
date: 2024-05-14T19:27:37+10:00
weight: 3
---

`ACDC` executables are released for Windows, MacOS, and Linux, which can be downloaded from the project's [Github Releases](https://github.com/openfast/acdc/releases). It may also be compiled from source using [WAILS](https://wails.io/) after following its [Installation Instructions](https://wails.io/docs/gettingstarted/installation). 

### Install OpenFAST

`ACDC` uses OpenFAST to run the turbine simulations and perform the linearization so a working installation of OpenFAST is needed. For detailed installation instructions, please refer to [https://openfast.readthedocs.io](https://openfast.readthedocs.io/en/main/source/install/index.html). Linearization works best when OpenFAST is compiled with the `DOUBLE_PRECISION=ON` option.

`ACDC` attempts to be independent of OpenFAST versions; however, the model files (OpenFAST input files) much match the version of OpenFAST selected in `ACDC`.

### Test OpenFAST Installation

Running the command `openfast` in a terminal will produce output similar to the following:

```
**************************************************************************************************
 OpenFAST

 Copyright (C) 2024 National Renewable Energy Laboratory
 Copyright (C) 2024 Envision Energy USA LTD

 This program is licensed under Apache License Version 2.0 and comes with ABSOLUTELY NO WARRANTY.
 See the "LICENSE" file distributed with this software for details.
 **************************************************************************************************

 OpenFAST-v3.5.3-dirty
 Compile Info:
  - Compiler: GCC version 13.2.0
  - Architecture: 64 bit
  - Precision: double
  - OpenMP: No
  - Date: Apr 17 2024
  - Time: 12:48:37
 Execution Info:
  - Date: 05/14/2024
  - Time: 11:03:46-0400


  Syntax is:

     OpenFAST [-h] <InputFile>

  where:

     -h generates this help message.
     <InputFile> is the name of the required primary input file.

  Note: values enclosed in square brackets [] are optional. Do not enter the brackets.

  Invalid syntax: no command-line arguments given.

  Aborting OpenFAST.
```

### Run ACDC

Start `ACDC` by running the executable downloaded from Github. `ACDC` will start and display the `Project` page as shown below (recent file paths have been removed):

![Project](project.png)

### Troubleshooting

#### macOS: developer cannot be verified

On macOS, users who have downloaded the app will usually see the error below the first time it is opened.
The error message states: 

```
"ACDC.app" cannot be opened because the developer cannot be verified.
macOS cannot verify that this app is free from malware.
```

{{< figure src="unsigned-error.png" width="350" >}}

This is due to the app distribution not being "signed" through Apple's developer certificate process.

To resolve the issue, right-click on the `ACDC.app` icon and select `Open`.
Then, click `Open` on the dialog that is displayed asking for verification to open an unsigned app.
Alternatively, the [extended attribute](https://en.wikipedia.org/wiki/Extended_file_attributes#macOS) can be removed with the [xattr](https://ss64.com/mac/xattr.html) command line tool:

```bash
xattr -d com.apple.quarantine ACDC.app
```
