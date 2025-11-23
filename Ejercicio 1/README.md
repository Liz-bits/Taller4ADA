## Lista de amistades bidireccionales:

1.  Alice - Bob
2.  Bob - Charlie
3.  Charlie - Drake
4.  Alice - Emma
5.  Emma - Frank
6.  Bob - Grace
7.  Grace - Henry
8.  Drake - Isabel
9.  Isabel - Jack
10. Frank - Jack
11. Alice - Henry
12. Charlie - Grace
13. Emma - Isabel

## Tabla de conexiones por usuario

| Usuario | Amigos directos (con separación 1) |
|---------|-------------------------------|
| **Alice** | Bob, Emma, Henry |
| **Bob** | Alice, Charlie, Grace |
| **Charlie** | Bob, Drake, Grace |
| **Drake** | Charlie, Isabel |
| **Emma** | Alice, Frank, Isabel |
| **Frank** | Emma, Jack |
| **Grace** | Bob, Charlie, Henry |
| **Henry** | Alice, Grace |
| **Isabel** | Drake, Emma, Jack |
| **Jack** | Frank, Isabel |

## Casos de prueba esperados

### Test 1: Alice,  N = 1

Conexiones directas de amistad con Alice

```
Usuario: Alice
N = 1
Conexion con: Bob, Emma, Henry
```

### Test 2: Alice,  N = 2

Los amigos de las conexiones de Alice (Amigos de Bob, Emma y Henry)

```
Usuario: Alice
N = 2
Conexion con: Charlie, Frank, Grace, Isabel
```

### Test 3: Alice,  N = 3

Los amigos de los amigos de las conexiones de Alice (sin repetirse en las conexiones pasadas)

```
Usuario: Alice
N = 3
Conexion con: Drake, Jack
```
