# PageRank-BasesDatosAvanzadas
Proyecto Bases de datos avanzadas GR40 - TEC
Tecnológico de Costa Rica
Profesor
Jose Castro
Estudiantes
	Cristian Arias Chaves
	Johel Rodriguez
	Oscar Blandino

##Instrucciones
El programa espera un archivo con nombre TestFile.txt ubicado en la carpeta Input. Este archivo tiene el siguiente formato:
*La primera línea corresponde al número de nodos del archivo
*A partir de la primera línea, cada número de línea - 2 corresponde al número de nodo, siendo el primer nodo el 0.
*Cada nodo puede o no puede tener aristas que salen a otros nodos, estas se representan con una lista de los nodos destinos separados por espacio.

Se puede generar un archivo de prueba ejecutando el script generador2.py, al ejecutarlo se le solicitara la cantidad de nodos que desea tenga el archivo, el generador automáticamente generará un mínimo de 0 y máximo de 100 nodos salientes para cada nodo.

Para ejecutar el programa ejecute el ejecutable llamado PageRank-BasesDatosAvanzadas.exe, la consola le solicitará un valor para betha (el cuál será usado como Damping al calcular el page rank para cada nodo).

##El programa

*El page rank inicial esta dado por 1 / cantidad de nodos.
*El pograma cuenta de dos ciclos de map/reduce, que se ejecutan iterativamente hasta que todos los elementos del vector de page ranks convergan con su antecesores con una diferencia de +/- 0.000000001.
#El primer ciclo map/reduce calcula el valor del page rank para cada arista entrante y realiza la sumatoria.
#El segundo ciclo map/reduce calcula el valor del page rank para cada nodo y verifica si todos los nodos convergieron a la diferencia deseada.
*El programa hace uso de rutinas y canales de go para emular un ambiente distribuido, por lo que los ciclos de cada map reduce no se ejecutan en un orden determinístico.

##Salida
El resultado del vector resultante se puede ver en el archivo de salida Results.txt, el cuál tiene el siguiente formato:
*El número de línea corresponde al número de nodo y el valor al page rank del nodo.
*La precisión con que se imprime en el archivo es de 25 decimales.




