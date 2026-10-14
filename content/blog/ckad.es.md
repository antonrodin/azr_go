---
title: Mi experiencia con el examen CKAD
slug: experiencia-ckad
description: Suspendí el CKAD en 5 minutos y lo aprobé a la segunda. Las reglas que nadie cuenta, el entorno del examen, la gestión del tiempo y cómo prepararlo en 60-90 horas.
date: 2026-10-14
image: /public/img/blog/ckad.jpg
image_alt: Logo de Kubernetes con el texto CKAD
image_credit: Logo de Kubernetes, dominio público (Wikimedia Commons)
---

Hace unos días me presenté por segunda vez al examen CKAD y lo he aprobado. Efectivamente, la primera vez no pude aprobar; de hecho, mi primer examen duró unos 5 minutos… De ahí que me gustaría contar un poco la experiencia del proceso, para que quizás no cometas mis mismos errores. He visto algún vídeo similar y, sinceramente, no cuentan algunas cosas importantes que hay que tener en cuenta cuando te vas a presentar a este examen.

Antes de nada, un poco de historia de por qué me he presentado y cuánto he tardado. En junio elegí la certificación CKAD para entender bien cómo funciona Kubernetes orientado a desarrollo. No quiero ser DevOps, pero quiero saber cómo funciona, ya que actualmente trabajo con microservicios que se despliegan con Kubernetes y Argo. Empecé aproximadamente a mediados o finales de junio, no me acuerdo muy bien. Le he ido dedicando entre 30 minutos y 1 hora al día. Diría que más bien 30 minutos: siendo padre no puedes sacar mucho más. Calculo que en total he dedicado entre 60 y 90 horas a prepararme el examen.

## Revisa las reglas muy bien

Para no haceros perder el tiempo, os doy unos consejos y os cuento por qué cateé la primera vez… Bueno, realmente no cateé: me echaron del examen a los 5 minutos. Una de las reglas del CKAD, del CKA y de cualquier certificación de The Linux Foundation es no hablar y **no mover los labios**. Es decir, si tenéis la extraña costumbre de leer los problemas o el enunciado, o simplemente de comentar el problema con vosotros mismos, tenéis que aprender a evitarlo. Está prohibido. El revisor puede considerar en ese momento que estás hablando con alguien y te echa del examen… y ahí se acaba. No es ni la primera ni la última vez que pasa (hay algún hilo de Reddit al respecto). En mi caso me quedé a cuadros, porque no entendí del todo por qué me echaban del examen en el primer ejercicio. El entorno es un pelín hostil la primera vez.

Después del fiasco, gracias a ChatGPT y a unos casos similares pude entender mejor la razón y reclamar. Es decir, escribí un ticket de reclamación y me contestaron. Por ser la primera vez, me devolvieron la posibilidad de hacer otro examen. Si no, os toca pasar por caja otra vez, lo que supone un desembolso de más de 250 dólares (y eso con oferta). Así que me pude presentar una segunda vez, y esta vez me preocupé de entrenar en killer.sh, que es lo más parecido que te vas a encontrar al examen real. Esa primera vez fue en septiembre; estaba bastante preparado y, aun así, mi examen duró 5 minutos. Os puedo asegurar que afecta bastante al estado de ánimo. Además, The Linux Foundation no se caracteriza por resolver tus problemas rápidamente: la respuesta tardó unos días, si no me equivoco. Me alejé una semana del CKAD y básicamente perdí un mes en volver a repasarlo todo y presentarme.

Lógicamente hay más reglas que te pueden afectar, así que cuidado. Te puede tocar un entorno muy, muy lento que provoque tu suspenso; hay quejas al respecto. Repasa bien las reglas e intenta cumplirlas, échale un ojo a los problemas que ha tenido la gente en Reddit y otras comunidades: es muy importante. Reclama si consideras que has tenido una experiencia injusta o insatisfactoria.

Por otro lado, el examen CKAD no es nada fácil, aun siendo introductorio; me ha recordado a los exámenes de la universidad. Son 2 horas y no te va a dar tiempo a acabarlo todo. Es decir, lo que hagas tiene que estar bien hecho, hay que tener cuidado con eso. Vas a ir con algo de nervios y te pueden jugar una mala pasada. Tómate la primera convocatoria como «introductoria»: no me extrañaría que, por cosas del entorno o por nervios, la suspendas. Ahora volvamos al entorno.

## El entorno no es el mismo que tienes en casa

Seguramente para trabajar tienes un ordenador con macOS o quizás Windows; normalmente los developers tenemos Mac. El entorno del examen es Linux (ni idea de qué distribución), por lo que muchos de tus comandos y atajos de teclado no te van a funcionar. Mi recomendación es hacer obligatoriamente, una semana antes del examen, un par de veces las sesiones de killer.sh. Te van a preparar más que cualquier curso que haya por ahí. Es lo más parecido que hay.

Lo que normalmente tardas 5 minutos en hacer en tu equipo, puedes tardar el doble en hacerlo en el entorno del examen, porque irá mucho más lento, no te resulta tan familiar… etc. Ante todo, calma y haz las cosas despacio en ese entorno. Tómatelo como otra sesión de killer.sh y apunta mentalmente las cosas que te preguntan para practicarlas después. Si apruebas, bien, y si no, lo harás la segunda vez. En principio tienes dos intentos, aunque cuidado: si te echan como a mí, te pueden quitar el segundo intento si consideran que has intentado hacer trampa. Con lo cual tendrás que pasar por caja una vez más.

En mi caso, aunque la segunda vez ya sabía de qué iba el tema, tuve unos problemas raros con vim. Era imposible guardar un archivo. Después, debatiéndolo con una IA, posiblemente el problema estaba en escribir los «:», con lo cual la combinación ESC + `:wq` no funcionaba siempre. A día de hoy sigo sin tener ni idea de por qué… pero aun así pude sacar el examen adelante.

Cuidado con el monitor. No se te ocurra hacer el examen en un portátil. Es casi obligatorio hacerlo con un monitor decente, lo que te obliga a comprar una webcam. No está permitido tener dos pantallas: tienes que tener solo un monitor. Hacer el examen en una pantalla pequeña es el doble de difícil. Avisados quedáis. Otro de mis errores fue presentarme con mi MacBook Air de 13 pulgadas… No era capaz ni de hacer un CTRL+F en el navegador, que creo que es Firefox por defecto.

Los revisores (*proctors*, creo que se llaman) no te van a solucionar mucho. Por regla general están ahí para vigilarte; puedes hablar con ellos por chat, pero no te van a ayudar en nada. Al final del examen intenté preguntar qué me pasaba con vim y simplemente me contestaron que sentían que la experiencia no estuviera siendo buena. En el primer examen, de hecho, la experiencia resultó muy, muy hostil, porque le estaba preguntando al revisor qué estaba haciendo mal y el resultado fue echarme del examen. Así que no esperéis mucha ayuda por ahí. Son bastante estrictos, hay muchas quejas al respecto por internet, pero es lo que hay.

## Gestión del tiempo

Practica la gestión del tiempo en killer.sh; para ello tienes la herramienta «flag». Sé que lo comenta todo el mundo, pero es imprescindible. No te quedes atascado en ningún ejercicio. Primero haz los fáciles, o los que consideres fáciles. Por fáciles entiende los que se pueden hacer en 5-10 minutos sin problema. Posiblemente serán un 30-40 %. Además, hacer esos fáciles te quitará los nervios iniciales.

Tendrás otros que son fáciles pero requieren más tiempo. Los raritos déjalos para el final. La regla aquí es no atascarse: intenta hacer primero lo fácil para tener algo seguro y después ve subiendo la nota con lo más complicado, poco a poco. Posiblemente a algunos no te dé tiempo ni a intentarlos. De hecho, siempre meten alguna pregunta que ni siquiera forma parte del temario, rollo Docker o Helm raros. De esas, casi seguro, olvídate… A mí me cayó una de Docker para crear una imagen y etiquetarla. Simplemente no me acordaba del comando para el tag. Pasé de la pregunta, ya que la dejé para lo último.

## No pierdas demasiado tiempo en cursos raros: ChatGPT es tu amigo

Puedes usar, por ejemplo, el famoso curso de Udemy, que te da acceso a KodeKloud para practicar. Personalmente no lo he visto muy necesario, aunque algo ayuda; tengo el curso, pero no lo acabé al 100 % porque además mete mucha paja innecesaria. Lo que sí ayuda es usar sabiamente una IA cualquiera: aliméntala con el contexto de todos los ejercicios que hagas y que te entrene para el examen. En mi caso tenía un proyecto en ChatGPT que me creaba ejercicios de repaso. La IA estaba alimentada con el curso oficial y con los ejercicios de killer.sh. Es decir, te crea ejercicios similares y encuentra tus puntos débiles: donde más fallabas, digamos, te creaba ejercicios para reforzarlo.

Para ahorrar algo de tiempo, en el curso oficial te puedes saltar toda la teoría, porque es una pérdida de tiempo y está mal redactada. Un churro infumable, digamos. En el curso de Udemy preocúpate solo de hacer los laboratorios y ver las soluciones. Ese curso sí lo puedes usar para la teoría, que es menos infumable. Por orden de prioridad:

1. Laboratorios oficiales. Imprescindibles. Pásalos por una IA (Gemini, ChatGPT o la que quieras) para que te los adapte a Minikube.
2. Killer.sh.
3. ChatGPT, para que te cree nuevos retos basados en los dos anteriores.
4. Udemy + KodeKloud.

Haz un presupuesto de horas. Yo calculo que con 60 horas es suficiente para ir medianamente bien al examen. Te presentas y, si no apruebas, añade unas 20 horas más o menos. Unas 30 horas al curso oficial, sobre todo a los laboratorios pasados por Minikube y ChatGPT. No te pongas a montar un clúster en AWS, que es una pérdida de tiempo. Unas 10-15 horas a killer.sh y otras 15 horas de repaso con ChatGPT o similares. Digamos que, con un poco de disciplina y 30 minutos-1 hora al día, lo tienes en dos meses fácilmente. Lo único: hay que ser constante.

## Comandos que tienes que saber al dedillo… todo con namespace

Hay muchos comandos básicos que tienes que saber al dedillo, pero recuerda además que, si se te olvidan, los puedes buscar tal cual en la documentación de Kubernetes:

```bash
alias kk=kubectl
kubectl run …
kubectl create -f …
kubectl replace -f … --force
# … etc.
```

Haz tu propia lista, pero te tienen que salir de forma natural.

## No es tan difícil

Aunque es mucho más difícil que un examen tipo test de cualquier otra certificación, sigue siendo asequible. No te desanimes si suspendes la primera vez, porque es muy normal. El entorno del examen es muy hostil y no siempre funciona bien. Considéralo experiencia y poco más. Aunque no apruebes, la experiencia que adquieras con Kubernetes es muy valiosa; al menos para mí ha sido así. De hecho, ahora mismo ni siquiera necesito la certificación como tal, porque el conocimiento ya lo tengo. Me ha ayudado bastante a entender mi trabajo. En los tiempos que corren solo necesito entender cómo funcionan las cosas; no necesito el papel, aunque siempre viene bien.

Incluso si suspendes dos veces, date tiempo y mírate cualquier otra cosa. Hay cientos de certificaciones o cursos muy valiosos y mucho más fáciles, como AWS Cloud Practitioner o AI Practitioner, que además te sacas en un par de semanas.
