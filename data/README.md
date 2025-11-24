# Sample Data Seeds

Each `.csv` file contains dummy rows for the corresponding PostgreSQL table.

## Usage

1. Ensure the target database has the tables already created and the `uuid-ossp` extension enabled (los servicios ya lo hacen al iniciar).
2. Importa un archivo usando `psql` o `COPY ... FROM STDIN`. Ejemplo rápido:
   ```bash
   \COPY users FROM 'data/users.csv' WITH (FORMAT csv, HEADER true);
   ```
   Ajusta el nombre de la tabla según corresponda (`route`, `workout`, `leaderboards`, etc.).

Si prefieres colocarlos directamente en el directorio de datos que usa tu clúster, mueve los `.csv` desde `data/` a esa ruta antes de ejecutarlos.
