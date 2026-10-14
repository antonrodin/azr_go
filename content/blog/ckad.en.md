---
title: My CKAD exam experience
slug: ckad-experience
description: I failed the CKAD in 5 minutes and passed on the second try. The rules nobody tells you about, the exam environment, time management and how to prepare in 60-90 hours.
date: 2026-10-14
image: /public/img/blog/ckad.jpg
image_alt: Kubernetes logo with the text CKAD
image_credit: Kubernetes logo, public domain (Wikimedia Commons)
---

A few days ago I sat the CKAD exam for the second time, and I passed. Yes, the first time I didn't make it; in fact, my first exam lasted about 5 minutes… That's why I'd like to share a bit of the experience, so that maybe you don't make the same mistakes I did. I've watched a few similar videos and, honestly, they leave out some important things you need to keep in mind before you sit this exam.

First, a bit of background on why I took it and how long it took me. In June I picked the CKAD certification to really understand how Kubernetes works from a developer's point of view. I don't want to be a DevOps engineer, but I want to know how it works, since I currently work with microservices that are deployed with Kubernetes and Argo. I started around mid or late June, I don't remember exactly. I spent between 30 minutes and 1 hour a day on it. Closer to 30 minutes, to be honest: as a father you can't squeeze out much more. All in all, I estimate I spent between 60 and 90 hours preparing for the exam.

## Read the rules very carefully

To save you some time, here are a few tips, and the story of why I failed the first time… Well, I didn't actually fail: I got kicked out of the exam after 5 minutes. One of the rules of the CKAD, the CKA and any The Linux Foundation certification is that you must not talk and **must not move your lips**. In other words, if you have the odd habit of reading the tasks out loud, or just talking a problem through to yourself, you need to learn to stop. It's forbidden. The proctor may decide you're talking to someone and kick you out… and that's the end of it. It's neither the first nor the last time this has happened (there are a few Reddit threads about it). In my case I was stunned, because I didn't really understand why I was being thrown out during the very first task. The environment is a little hostile the first time.

After the fiasco, thanks to ChatGPT and a few similar cases, I understood the reason better and was able to file a complaint. That is, I opened a support ticket and they replied. Since it was my first time, they gave me back the chance to take another exam. Otherwise you'd have to pay again, which means over 250 dollars (and that's with a discount). So I was able to sit it a second time, and this time I made sure to train on killer.sh, which is the closest thing you'll find to the real exam. That first attempt was in September; I was well prepared and, even so, my exam lasted 5 minutes. I can assure you it takes a toll on your morale. On top of that, The Linux Foundation isn't exactly known for solving your problems quickly: the reply took a few days, if I remember correctly. I stepped away from the CKAD for a week and basically lost a month going over everything again before retaking it.

Obviously there are more rules that can affect you, so be careful. You might get a very, very slow environment that makes you fail; people complain about it. Review the rules carefully and try to follow them, and take a look at the problems people have had on Reddit and other communities: it really matters. Complain if you feel your experience was unfair or unsatisfactory.

On the other hand, the CKAD exam is not easy at all, even though it's an entry-level one; it reminded me of university exams. It's 2 hours and you won't have time to finish everything. That means whatever you do has to be done right, so be careful about that. You'll be a bit nervous, and nerves can play tricks on you. Treat your first attempt as an "introductory" one: I wouldn't be surprised if you failed it because of the environment or nerves. Now, back to the environment.

## The environment is not the one you have at home

You probably work on a computer with macOS, or maybe Windows; most of us developers use a Mac. The exam environment is Linux (no idea which distribution), so many of your usual commands and keyboard shortcuts won't work. My recommendation: in the week before the exam, you absolutely must go through the killer.sh sessions a couple of times. They'll prepare you better than any course out there. It's the closest thing there is.

Something that normally takes you 5 minutes on your own machine can take twice as long in the exam environment, because it's much slower, it's less familiar… and so on. Above all, stay calm and take it slowly in that environment. Treat it as one more killer.sh session and make a mental note of what they ask so you can practise it later. If you pass, great; if not, you'll pass the second time. In principle you get two attempts, but careful: if they kick you out like they did with me, they can take away your second attempt if they think you tried to cheat. Then you'll have to pay all over again.

In my case, even though the second time I knew what I was in for, I had some weird problems with vim. It was impossible to save a file. Later, discussing it with an AI, it seemed the problem was probably typing the ":" character, so the ESC + `:wq` combo didn't always work. To this day I have no idea why… but I still managed to get through the exam.

Watch out for the monitor. Don't even think about taking the exam on a laptop. It's almost mandatory to take it with a decent monitor, which means you'll have to buy a webcam. You're not allowed to have two screens: only one monitor. Taking the exam on a small screen is twice as hard. You've been warned. Another of my mistakes was sitting it on my 13-inch MacBook Air… I couldn't even manage a CTRL+F in the browser, which I believe is Firefox by default.

The proctors won't fix much for you. As a rule they're there to watch you; you can talk to them through the chat, but they won't help you with anything. At the end of the exam I tried asking what was going on with vim, and they simply replied that they were sorry the experience wasn't a good one. In the first exam the experience was actually very, very hostile: I was asking the proctor what I was doing wrong, and the result was getting kicked out. So don't expect much help there. They're quite strict, there are plenty of complaints about it online, but that's how it is.

## Time management

Practise time management on killer.sh; that's what the "flag" tool is for. I know everybody says it, but it's essential. Don't get stuck on any task. Do the easy ones first, or the ones you consider easy. By easy I mean the ones you can do in 5-10 minutes without breaking a sweat. They'll probably be around 30-40 %. Plus, doing the easy ones will get rid of the initial nerves.

There will be others that are easy but take longer. Leave the weird ones for the end. The rule here is not to get stuck: do the easy stuff first to bank some points, and then raise your score bit by bit with the harder ones. You probably won't even have time to try some of them. In fact, they always throw in a question that isn't even part of the syllabus, some odd Docker or Helm thing. Forget about those, almost certainly… I got a Docker one about building an image and tagging it. I just couldn't remember the tag command. I skipped it, since I'd left it for last.

## Don't waste too much time on random courses: ChatGPT is your friend

You can use, for example, the well-known Udemy course, which gives you access to KodeKloud for practice. Personally I didn't find it very necessary, although it helps a bit; I have the course, but I didn't finish it 100 % because it also includes a lot of unnecessary filler. What does help is using any AI wisely: feed it the context of every exercise you do and let it train you for the exam. In my case I had a project in ChatGPT that created review exercises for me. The AI was fed with the official course and the killer.sh exercises. In other words, it creates similar exercises and finds your weak spots: wherever you failed most, it would create exercises to fix that.

To save some time, you can skip all the theory in the official course, because it's a waste of time and badly written. An unreadable mess, let's say. In the Udemy course, focus only on doing the labs and checking the solutions. You can use that course for the theory, which is less painful. In order of priority:

1. Official labs. Essential. Run them through an AI (Gemini, ChatGPT or whatever you like) so it adapts them to Minikube.
2. Killer.sh.
3. ChatGPT, to create new challenges based on the two above.
4. Udemy + KodeKloud.

Make an hours budget. I reckon 60 hours is enough to go into the exam reasonably well prepared. You sit it and, if you don't pass, add roughly another 20 hours. About 30 hours on the official course, mostly the labs run through Minikube and ChatGPT. Don't bother setting up a cluster on AWS, it's a waste of time. About 10-15 hours on killer.sh and another 15 hours of review with ChatGPT or similar. Let's say that, with a bit of discipline and 30 minutes to 1 hour a day, you'll easily have it done in two months. The only thing: you have to be consistent.

## Commands you need to know by heart… all with namespaces

There are many basic commands you need to know by heart, but also remember that if you forget them, you can look them up as-is in the Kubernetes documentation:

```bash
alias kk=kubectl
kubectl run …
kubectl create -f …
kubectl replace -f … --force
# … etc.
```

Make your own list, but they have to come naturally.

## It's not that hard

Although it's much harder than a multiple-choice exam for any other certification, it's still doable. Don't get discouraged if you fail the first time, because it's very common. The exam environment is very hostile and doesn't always work well. Consider it experience and little more. Even if you don't pass, the experience you gain with Kubernetes is very valuable; at least it has been for me. In fact, right now I don't even need the certification as such, because I already have the knowledge. It has helped me a lot to understand my job. These days I just need to understand how things work; I don't need the paper, although it always comes in handy.

Even if you fail twice, give it time and look into something else. There are hundreds of very valuable and much easier certifications or courses, like AWS Cloud Practitioner or AI Practitioner, which you can also get in a couple of weeks.
