> **Warning**: this project is still under heavy development and should not be considered ready for usage!
> Both the api and underlying implementation are subject to change
# Muxillary
It's the dumbest pun I could think of.

Muxillary is a simple http route parser and multiplexer.

It aims to integrate seamlessly with the golang http package, while providing functionality missing in the standard library.

# Quickstart
## Installation 
`go get github.com/Sibo9436/muxillary`

## Example
```go
import (
    "net/http"
    "github.com/Sibo9436/muxillary"
)


func main(){
    mux := muxillary.NewMuxillaryHandler("/")
    mux.Get("test/:id", func (w http.ResponseWriter, r *http.Request){
        id := muxillary.Value("id", r)
        //Do something with your request
        //-----
    })
    http.ListenAndServe(":3000", mux)
}
```

# Ideas
Secondo me potrebbe essere divertente implementarlo utilizzando una struttura ad albero:

```mermaid
graph TD;
    root-->A(/cavalli);
    root-->C(/proprietari);
    A-->B(/proprietari)
    C-->x(/:id)
    B-->y(/:id)
    A-->z(/:id)
```

Bisogna però vedere che effetto potrebbe avere sulla performance

In questo modo, quando vado a popolare i miei endpoints se trovo un endpoint che potrebbe creare un conflitto posso andare in errore gioiosamente

> Suggestions welcome: should muxillary have a default way of handling errors? Or should we keep it as close to the http.Handler api as possible?

## Rules 
Per il mapping mi ispiro almeno parzialmente a SpringBoot

Intanto devo dare la precedenza ai path definiti:

Ad esempio se ho `/contatti` e `/:nome` prima controllo che la parola non sia contatti 
e poi passo a nome 

In secondo luogo non devono essere permesse due path variables allo stesso livello
ovvero `/:first` `/:second` non è una descrizione corretta, mentre
`:first` `:second/count` lo è. Questo è forse il problema più interessante da risolvere
ovvero capire come non usare regexp ma comunque essere in grado di matchare correttamente il path inserito

Inoltre lo standard openAPI definisce le pathvariables con {nome} e usa ; e * per indicare che formato 
ci si aspetta di ricevere per oggetti come mappe o liste 

Per adesso questa funzionalità resta OOS ma non è detto che non venga aggiunta later.

Un'idea iniziale potrebbe essere quella di avere un nodo con una mappa di possibili figli e un singolo nodo figlio "any"
che a sua volta avrà una sua mappa di possibili figli e una mappa di possibili valori per il nome del parametro
IN REALTA', sarà una tamarrata ok, ma si potrebbe semplicemente popolare tutti i valori possibili con il valore rilevato e risolvere il 
problema del lookeahed

Ad esempio se il la mia pathvaribale a quell'altezza può essere sia :first che :second, valorizzo semplicemente entrambe


## Todo 
Per adesso siamo riusciti a superare un test molto basilare
Bisogna andare a vedere un pochino come ci aspettiamo essere una specifica decente per come gestire i path parameters 

Oltre ai path parameters devo anche trovare un modo decente per parsare i query parameters e schiantarli in una mappa, perché mi sa che golang
di suo non lo fa

La cosa più importante sarà poi fare dei benchmark e capire se e come modificare il sistema dei path per evitare problemi di scaling

Per la questione performance allo stato attuale non so bene quanto pesi il fatto di avere degli oggetti interi per il path

Importante sarà poi gestire le collisioni tra path!!

Nel frattempo mi stanno venendo in mente altre diciottomila idee per creare un sistema più simile a quello ad esempio di SpringBoot 
ma in realtà non ha molto senso visto che l'idea principale è mantenere un mux che sia perfettamente compatibile con il pacchetto http
e che sia il meno ingombrante possibile

Tra le cose da fare devo poter mettere come minimo la possibilità di definire una funzione che restituisca un 404 customizzato, non dovrebbe essere difficile



