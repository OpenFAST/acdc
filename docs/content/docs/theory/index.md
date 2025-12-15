---
title: 'Theory'
date: 2024-05-14T19:27:37+10:00
weight: 10
---

## Overview

A campbell diagram is a powerful visualization tool used in wind turbine analysis to understand how natural frequencies and damping ratios vary with rotor and wind speeds. It helps identify potential resonance conditions and instabilities where structural modes may be excited by harmonic frequencies, which is critical for avoiding fatigue damage and ensuring safe operation.

The campbell diagram generation process with OpenFAST follows a systematic automated workflow consisting of five main steps:

1. **Import OpenFAST Model** - Analyze the turbine to get time domain and linearization data
2. **Define Operating Points** - Specify rotor speeds and other operating conditions for linearization
3. **Run Simulations** - Execute OpenFAST linearization at each operating point
4. **Multi-Blade Coordinate (MBC) Transformation** - Transform rotating blade coordinates to fixed reference frame
5. **Eigenanalysis and Mode Identification** - Extract natural frequencies and mode shapes, then track modes across operating points

![](workflow.png)

## Multi-Blade Coordinate (MBC) Transformation

For wind turbines with multiple blades, the dynamics in the rotating reference frame are periodic and coupled due to rotor rotation. The Multi-Blade Coordinate transformation converts the rotating blade degrees of freedom into a fixed (non-rotating) reference frame, producing:

- **Collective modes** - All blades move in phase (symmetric motion)
- **Cyclic modes** - Blades move with phase differences (asymmetric motion)
- **Fixed-frame representation** - Time-invariant system matrices suitable for eigenanalysis

This transformation is essential because it:
- Decouples the periodic equations of motion into time-invariant form
- Enables conventional eigenanalysis techniques
- Separates symmetric and asymmetric rotor modes
- Simplifies the identification of structural modes

The MBC transformation yields state-space matrices in the fixed reference frame, which can then be analyzed using standard linear algebra techniques.

## Eigenanalysis

Eigenanalysis is a mathematical technique used to compute modal and stability characteristics, decomposing a linear system into its fundamental components: eigenvalues and eigenvectors.

- **Eigenvalues** - Complex numbers representing natural frequencies and damping ratios of the system modes
- **Eigenvectors** - Complex vectors describing the mode shapes, indicating how different degrees of freedom participate in each mode

For wind turbine analysis, eigenanalysis of the MBC-transformed system provides:
- Natural frequencies at each operating point (used to construct the Campbell diagram)
- Mode shapes that characterize the physical motion (e.g., tower bending, blade flap/edge, drivetrain torsion)
- Damping characteristics indicating stability of each mode

The eigenvalues are typically plotted against rotor speed to create the Campbell diagram, with rotational harmonic lines overlaid to identify potential resonances.

## Mode Identification

Mode identification is the process of tracking the same physical mode across different operating points based on the eigenvectors and eigenvalues produced by MBC transformation and eigenanalysis. This is a challenging task because:

- Natural frequencies of different modes may cross or veer as operating conditions change
- Mode shapes may gradually evolve with rotor speed
- Multiple modes with similar characteristics may be present
- Numerical noise can affect mode ordering

As seen from the figure below, some modes can be clearly distinguished across the operating range, whereas others may cross or come very close to each other, making automated tracking difficult.

![](mode_identification.png "Modes for NREL 5MW")

### Similarity Metrics

Modal identification relies on quantitative similarity measurements to determine which modes at different operating points correspond to the same physical phenomenon.

**Modal Assurance Criterion (MAC)**

The Modal Assurance Criterion compares the complex eigenvectors of two modes to quantify their similarity. The MAC value ranges from 0 (completely dissimilar) to 1 (identical mode shapes).

<!-- {{< figure src="mac.png" width="300" >}} -->

$$
\text{MAC}(\mu_1, \mu_2) = \left( \frac{|\mu_1^*\ \mu_2|}{||\mu_1||\ ||\mu_2||} \right)^2
$$

where \(\mu_1\) and \(\mu_2\) are complex eigenvectors from two different operating points, and \(^*\) denotes the complex conjugate transpose.

**Pole-Weighted MAC (MACXP)**

MACXP enhances the standard MAC by incorporating both eigenvector similarity and eigenvalue proximity:

<!-- {{< figure src="macxp.png" width="500" >}} -->

$$
\text{MACXP}(\mu_1, \mu_2) = \frac{\left(\frac{|\mu_1^*\ \mu_2|}{|\overline{\lambda_1} + \lambda_2|} + \frac{|\mu_1^{\top}\ \mu_2|}{|\lambda_1 + \lambda_2|}\right)^2}{\left(\frac{\mu_1^*\ \mu_1}{2|\text{Re } \lambda_1|} + \frac{|\mu_1^{\top}\ \mu_1|}{2|\lambda_1|}\right) \left(\frac{\mu_2^*\ \mu_2}{2|\text{Re } \lambda_2|} + \frac{|\mu_2^{\top}\ \mu_2|}{2|\lambda_2|}\right)}
$$

where \(\lambda_A\) and \(\lambda_B\) are the eigenvalues corresponding to modes A and B. The exponential weighting factor penalizes modes with dissimilar eigenvalues, providing a more robust similarity measure that considers both shape and frequency content.

