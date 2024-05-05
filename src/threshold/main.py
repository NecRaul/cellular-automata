import numpy as np
import matplotlib.pyplot as plt

size = np.linspace(1, 1000, 200)
iterations = np.linspace(1, 500, 200)

X, Y = np.meshgrid(size, iterations)
Z = np.log2(np.sqrt(X * Y))

plt.contourf(X, Y, Z, cmap="plasma")
plt.xlabel("Size", fontsize=16)
plt.ylabel("Iterations", fontsize=16)
plt.title("Threshold", fontsize=18)
plt.tight_layout()
plt.colorbar()
plt.grid()
plt.show()
