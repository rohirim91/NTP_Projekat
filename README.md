#
## Napredne Tehnike Programiranja – Predlog projekta

# Segmentacija slike upotrebom PSO algoritma

# Opis problema i pristup rešavanju

Ideja projekta jeste implementacija aplikacije koja bi omogućila segmentaciju gray-scale slika proizvoljnih dimenzija.

Tehnika segmentacije bi bila multithresholding, odnosno slika bi se delila u više segmenata upotrebom pragova. Primer sa 2 praga i tri segmenta:

Prag 1 – RGB(128,128,128), svi pikseli slike sa RGB vrednošću nižom od 128,128,128 će biti u grupisani u jedan segment.

Prag 2 – RGB(197,197,197), svi pikseli sa RGB vrednošću između 128,128,128 (prethodni prag) i 197,197,197 biće grupisani u ovaj prag.

Svi pikseli sa RGB vrednošću iznad praga 2 biće grupisani u treći segment.

Kako se pragovi ne bi birali ručno, planirani sistem će ovo automatizovati upotrebom Tsallis entropije (funkcije) koja se, uz jednostavniju i manje efikasnu Otsu funkciju, najčešće koristi za rešavanje problema određivanja optimalnih pragova. Ukratko, veća vrednost ove funkcije u teoriji znači i bolje pragove kojima se slika segmentiše.

Uzimajući gorenavedeno u obzir, problem segmentacije slike postaje optimizacioni problem. Kako bi se vrednost Tsallis funkcije maksimizovala, u sistemu će kao optimizacioni algoritam biti korišćen [PSO (particle swarm optimization) algoritam](https://en.wikipedia.org/wiki/Particle_swarm_optimization).

Ukratko, PSO algoritam je iterativni algoritam koji se zasniva na simulaciji skupa čestica koje se kreću kroz n-dimenzioni prostor pretrage i izračunavaju vrednost funkcije u svom položaju. Položaj čestice predstavlja koordinate u n dimenzija koje predstavljaju ulazne vrednosti funkcije koja se optimizuje.

Dodatno, korisniku će biti omogućeno da sliku pre segmentacije obradi [Superpixels tehnikom](https://arxiv.org/pdf/1612.01601.pdf) segmentacije. Suština ove tehnike leži u tome što se pikseli slike sa bliskom RGB vrednošću grupišu u celine – superpixele sa uprosečenim RGB vrednostima. Intuicija za korišćenje ove tehnike jeste da bi se segmentisanjem slike pomoću superpixel-a veoma malih dimenzija (red veličine 4-6 piksela) postigao efekat uklanjanja sitnih detalja koji korisniku nisu od značaja prilikom segmentisanja određene slike, a koji utiču na kvalitet segmentacije i to bez uklanjanja krupnijih detalja koji su korisniku relevantni. Takođe, omogućava se da se u daljem postupku segmentacije ne posmatraju pikseli, već superpikseli, što je od velikog značaja prilikom segmentacije slika velikih dimenzija (umesto, recimo 512\*512 piksela, obrađivalo bi se oko 200-300 superpiksela, zavisno od njihove dimenzije).

# Arhitektura sistema

Sistem bi bio sastavljen od 2 mikroservisa implementirana u Go jeziku i desktop klijentske aplikacije implementirane u Pharo jeziku:

- Pharo klijent – omogućava unos lokacije slike u fajl sistemu, unos željene lokacije na kojoj će se čuvati dobijene slike, unos parametara superpixels i PSO algoritma i vizuelizacija prostora pretrage ukoliko je slika segmentisana sa jednim ili dva praga.
- Mikroservis 1 – implementacija superpixels algoritma, učitavanje slike, njena obrada i čuvanje dobijene slike.
- Mikroservis 2 – implementacija segmentacije upotrebom PSO algoritma i Tsallis funkcije, učitavanje slike, čuvanje segmentisane slike i podataka neophodnih za vizuelizaciju prostora pretrage.

Komunikacija bi se odvijala putem Representational State Transfer (REST).

Kao što je već rečeno, Pharo klijent bi korisniku prikazao i pretraženi prostor pretrage PSO algoritma u slučajevima korišćenja jednog i dva praga.

Vizuelizovali bi se svi položaji čestice od prve do poslednje iteracije, budući da je ovakva vizuelizacija i najkorisnija korisniku, jer nudi prikaz celokupnog pretraženog prostora pretrage, što mu može biti od koristi ukoliko želi da ručno isproba neke kombinacije vrednosti pragova koje PSO algoritam nije pokrio. Primer su lokalni optimumi koji možda nisu dovoljno dobro istraženi, te jedan od njih može biti ustvari globalni optimum.

Vizuelizacija upotrebom jednog praga (binarna segmentacija), bi bila u obliku grafika funkcije, gde bi X koordinata predstavljala vrednosti praga u opsegu 0-255, dok bi Y koordinata predstavljala vrednost Tsallis funkcije.

Vizuelizacija upotrebom 2 praga bila bi u stilu heat mape, gde bi koordinate tačaka predstavljale vrednosti pragova, a boja tačke bi bila indikator vrednosti Tsallis funkcije u toj tački. Na primer, nijanse plave za niže vrednosti, a nijanse crvene za više (gradijent od plave ka crvenoj).

Primer vizuelizacije prostora pretrage:

![](RackMultipart20210319-4-gif1qy_html_68388fe14baab03a.png)

_Slika 1 - Vizuelizacija prostora pretrage u slučaju segmentacije sa jednim pragom (binarna segmentacija)_

![](RackMultipart20210319-4-gif1qy_html_5c582a2cbdb2b77d.png)

_Slika 2 - Vizuelizacija prostora pretrage u slučaju korišćenja dva praga_

# Dalja unapređenja

Neke ideje za dalja unapređenja sistema su:

- Paralelizacija superpixels algoritma
- Paralelizacija PSO algoritma
- Vizuelizacija po iteracijama PSO algoritma.
- Segmentacija slika u boji.
