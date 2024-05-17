# Changelog

All notable changes to this project will be documented in this file.

## v0.5.0-alpha

### Added

- Added [About](https://OpenFAST.github.io/acdc/docs/about) page to documentation
- Enabled [Zenodo](https://zenodo.org/) for repository to generate DOI for each release

## v0.4.0-alpha

### Added

- Documentation in Github Pages
- Apache 2 License

### Changed

- Moved repository to the [OpenFAST Organization](https://github.com/OpenFAST/acdc)

### Fixed

- Model: "Model" tab doesn't work [#6](https://github.com/deslaughter/acdc-app/issues/6)

## v0.3.0-alpha

### Added

- New example project based on NREL 5MW land based turbine (`examples/NREL_5MW_Land`)
- New example project based on NREL 5MW monopile turbine with HydroDyn and SubDyn (`examples/NREL_5MW_Monopile`)
- Results: Ability to load existing results and Campbell diagram when selecting linearization folder
- Analysis: Add Copy feature for Analysis cases [#3](https://github.com/deslaughter/acdc-app/issues/3)
- Results: Add button to visualize all operating points in line

### Changed

- Evaluate: SttsTime sets itself to 6, ignoring user input [#4](https://github.com/deslaughter/acdc-app/issues/4)
- Evaluate: Output defaults for OpenFAST [#5](https://github.com/deslaughter/acdc-app/issues/5)

### Fixed

- Results: Screen going blank when swapping diagram points
- Results: Results tab breaking when visualizing last OP [#2](https://github.com/deslaughter/acdc-app/issues/2)
- Model: "Model" tab doesn't work [#6](https://github.com/deslaughter/acdc-app/issues/6)
- Model: SubDyn filename in generated .fst-file and name of SubDyn-file in folder not compatible [#7](https://github.com/deslaughter/acdc-app/issues/7)

## v0.2.0-alpha

This release adds many features to ACDC

### Added

- Ability to customize Campbell Diagram
    - Line labels
    - Line colors
    - Swap points between lines
    - Filter non-structural modes (experimental)
- Visualize mode shapes for top, front, right-side, and isometric views of turbines
- Visualize path of nodes in mode shape animation to aid in mode identification

## v0.1.0-alpha

Initial alpha release of ACDC. This version has the following features:

- Create project
- Import OpenFAST model files
- Set recommended model defaults for performing linearization
- Configure analysis cases with wind speed, rotor speed, and blade pitch curve
- Generate operating points using spline interpolation of condition curve
- Run simulations and perform linearization at all operating points in parallel
- Import linearization files and perform multi-blade coordinate transform (MBC)
- Perform Eigenanalysis of transformed state matrix to get turbine natural frequency and damping
- Use modal assurance criteria (MAC) to connect modes across operating points
- Generate Campbell Diagram plot and display it to user
- Use spectral clustering to correct modes that transform across operating points