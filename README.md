# PageRank-BasesDatosAvanzadas
Proyecto Bases de datos avanzadas GR40 - TEC
Tecnol�gico de Costa Rica
Profesor
Jose Castro
Estudiantes
	Cristian Arias Chaves
	Johel Rodriguez
	Oscar Blandino

##Instrucciones
El programa espera un archivo con nombre TestFile.txt ubicado en la carpeta Input. Este archivo tiene el siguiente formato:
*La primera l�nea corresponde al n�mero de nodos del archivo
*A partir de la primera l�nea, cada n�mero de l�nea - 2 corresponde al n�mero de nodo, siendo el primer nodo el 0.
*Cada nodo puede o no puede tener aristas que salen a otros nodos, estas se representan con una lista de los nodos destinos separados por espacio.

Se puede generar un archivo de prueba ejecutando el script generador2.py, al ejecutarlo se le solicitara la cantidad de nodos que desea tenga el archivo, el generador autom�ticamente generar� un m�nimo de 0 y m�ximo de 100 nodos salientes para cada nodo.

Para ejecutar el programa ejecute el ejecutable llamado PageRank-BasesDatosAvanzadas.exe, la consola le solicitar� un valor para betha (el cu�l ser� usado como Damping al calcular el page rank para cada nodo).

##El programa

*El page rank inicial esta dado por 1 / cantidad de nodos.
*El pograma cuenta de dos ciclos de map/reduce, que se ejecutan iterativamente hasta que todos los elementos del vector de page ranks convergan con su antecesores con una diferencia de +/- 0.000000001.
#El primer ciclo map/reduce calcula el valor del page rank para cada arista entrante y realiza la sumatoria.
#El segundo ciclo map/reduce calcula el valor del page rank para cada nodo y verifica si todos los nodos convergieron a la diferencia deseada.
*El programa hace uso de rutinas y canales de go para emular un ambiente distribuido, por lo que los ciclos de cada map reduce no se ejecutan en un orden determin�stico.

##Salida
El resultado del vector resultante se puede ver en el archivo de salida Results.txt, el cu�l tiene el siguiente formato:
*El n�mero de l�nea corresponde al n�mero de nodo y el valor al page rank del nodo.
*La precisi�n con que se imprime en el archivo es de 25 decimales.




