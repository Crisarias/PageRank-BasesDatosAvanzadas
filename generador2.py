import os

import random


CantNodos = int(input("Cantidad de nodos "))

file = open("./filename.txt", "w")
#CantNodos = random.randrange(100)
print(CantNodos)
<<<<<<< Updated upstream
=======
file.write(str(CantNodos))
file.write(os.linesep)
>>>>>>> Stashed changes
#file.write(str(CantNodos)+os.linesep)
for lineas in range(CantNodos):
    cantvertex =  random.randrange(CantNodos)
    
 #   print(lineas,end = '')
    print(" ",end = '')
<<<<<<< Updated upstream
    file.write(str(lineas) + " ")
=======
    #file.write(str(lineas) + " ")
>>>>>>> Stashed changes
    lista = []
    for x in range(cantvertex):
        linkto= random.randrange(CantNodos)
        if linkto  not in lista:
            lista.append(linkto)
            file.write(str(linkto)+ " ")
            print(linkto,end = '')
            print(" ",end = '')
    print("")
    file.write(os.linesep)

file.close()