# Pokedex

Este proyecto es una implementación de una Pokedex en Go que interactúa con la API de PokeAPI para obtener información sobre Pokémon y áreas de ubicación. 

## Instalación

1. Clona el repositorio:
    ```sh
    git clone https://github.com/glopez94/pokedex.git
    ```
2. Navega al directorio del proyecto:
    ```sh
    cd pokedex
    ```
3. Instala las dependencias:
    ```sh
    go mod tidy
    ```

## Uso

Para ejecutar la Pokedex, utiliza el siguiente comando:
```sh
go run main.go
```

Una vez ejecutado, puedes utilizar los siguientes comandos dentro de la interfaz de línea de comandos:

- `exit`: Salir de la Pokedex.
- `help`: Muestra un mensaje de ayuda.
- `map`: Muestra las siguientes 20 áreas de ubicación.
- `mapb`: Muestra las 20 áreas de ubicación anteriores.
- `explore <area>`: Explora una área de ubicación específica.
- `catch <pokemon>`: Atrapa un Pokémon específico.
- `inspect <pokemon>`: Inspecciona un Pokémon capturado.
- `pokedex`: Lista todos los Pokémon capturados.

## Ejemplos

### Explorar un área de ubicación
```sh
Pokedex > explore kanto-route-1
Found Pokemon:
 - pidgey
 - rattata
```

### Atrapar un Pokémon
```sh
Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!
```

### Inspeccionar un Pokémon
```sh
Pokedex > inspect pikachu
Name: pikachu
Height: 4
Weight: 60
Stats:
  - speed: 90
  - special-defense: 50
  - special-attack: 50
  - defense: 40
  - attack: 55
  - hp: 35
Types:
  - electric
```

### Listar Pokémon capturados
```sh
Pokedex > pokedex
Your Pokedex:
 - pikachu
 - bulbasaur
```

## Estructura del Proyecto

- `main.go`: Contiene la lógica principal de la Pokedex y la interfaz de línea de comandos.
- `internal/pokecache/cache.go`: Implementa una caché simple para almacenar respuestas de la API.

## Contribuciones

Las contribuciones son bienvenidas. Por favor, abre un issue o envía un pull request.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT. Consulta el archivo `LICENSE` para más detalles.