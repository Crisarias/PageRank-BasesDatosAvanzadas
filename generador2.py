import os

import random


CantNodos = int(input("Cantidad de nodos "))

file = open("./inputFile/testFile.txt", "w")
#CantNodos = random.randrange(100)
print(CantNodos)
file.write(str(CantNodos)+"\n")
#file.write(os.linesep)
#file.write(str(CantNodos)+os.linesep)
for lineas in range(CantNodos):
    cantvertex =  random.randrange(CantNodos) - 1 
    if cantvertex > 100:
        cantvertex = 100
 #   print(lineas,end = '')
    print(lineas+1, end='')
    #file.write(str(lineas) + " ")
    lista = [lineas+1]
    for x in range(cantvertex):
        linkto = random.randrange(CantNodos)
        while linkto in lista:
            linkto= random.randrange(CantNodos)
       
        lista.append(linkto)
        file.write(str(linkto))
        if x != cantvertex-1:
            file.write(" ")
                        
    print("")
    file.write("\n")
#file.write(os.linesep)
file.close()
