# Payout Calculation Service

Servis PostgreSQL, Redis va HTTP API bilan ko'tarildi. API dokumentatsiyasi swagger orqali ochiladi:

```text
http://localhost:8080/swagger/index.html
```

## Ishga tushirish

```bash
docker compose up --build
```

Servis `http://localhost:8080` da ishlaydi

docker compose ichida `app`, `postgres`, `redis` va migration containerlari ko'tariladi, migrationlar app ishga tushishidan oldin ko'tariladi

default admin:

```text
username: admin
password: admin123
```

## Auth

`admin` barcha ma'lumotlarni ko'ra oladi, kuryer va order yaratadi, payout hisoblaydi, oylik cron jobni ishga tushira oladi

`courier` faqat o'ziga tegishli ma'lumotlarni ko'ra oladi.

## Payout hisoblash

Payout `courier_id` va oy bo'yicha hisoblanadi. Oy formati:

```text
YYYY-MM
```

misol uchun ->

```json
{
  "courier_id": "uuid",
  "period": "2026-09"
}
```

Hisobga faqat shu oy ichida `delivered` bo'lgan statusdegi orderlar kiradi

Komissiya bosqichlari:

```text
0-50 ta delivered order      10%
51-150 ta delivered order    12%
151+ ta delivered order      15%
```

Oy bo'yicha umumiy delivered statusdegi orderlar soniga qarab bitta stavka tanlanadi va shu stavka butun oylik summaga qo'llaniladi ->

```text
commission_amount = gross_amount * commission_rate
net_amount = commission_amount
```

Bu model tanlanganiga sabab bitta oy uchun bitta stavka bo'lsa, payoutni tekshirish va tushuntirish osonroq bo'ladi. Agar keyin order statusi o'zgarib, kuryer boshqa bosqichga o'tib qolsa, eski payout o'zgartirilmaydi. Oy qayta hisoblanadi va farq `payout_adjustments` jadvaliga yoziladi.

Kuryer oy o'rtasida ish boshlagan bo'lsa ham alohida prorate qilinmaydi. Hisob faqat shu kuryerga tegishli va shu oyda delivered bo'lgan orderlarga qaraladi

## Idempotency va concurrency

Bitta kuryer uchun bitta oyga faqat bitta payout bo'lishi mumkin.

```sql
CONSTRAINT uq_payouts_courier_period UNIQUE (courier_id, period)
```

Shu sababli bir xil `courier_id` va `period` bilan qayta calculate qilinsa, servis mavjud payoutni body ichida qaytaradi va status `409 Conflict` bo'ladi.

Parallel so'rovlar ham shu DB constraint orqali himoyalangan. Ya'ni 20 ta so'rov bir vaqtda kelsa, faqat bittasi insert qila oladi. Qolganlari unique constraint sabab mavjud payoutni qaytaradi.

## Adjustmentlar

Payout hisoblangandan keyin payout yozuvi o'chirilmaydi va ustidan yozilmaydi.

Agar hisoblangan oyga tegishli order statusi keyin o'zgarsa, servis bitta transaction ichida:

Isolation level ->

```text
1. payout rowni FOR UPDATE bilan lock qiladi
2. oy bo'yicha delivered count va gross amountni qayta hisoblaydi
3. oldingi adjustmentlarni hisobga oladi
4. faqat farqni payout_adjustments jadvaliga yozadi
```

`GET /payouts/{id}` javobida payout, adjustmentlar va yakuniy summa qaytadi:

```text
total_amount = payout.net_amount + sum(adjustments.net_delta)
```

## Oylik job

Har oyning 1-kunida o'tgan oy uchun barcha active kuryerlarga payout hisoblanadi.

Standart cron:

```text
PAYOUT_JOB_CRON=0 0 1 * *
```

Servis bir nechta instance bo'lib ishlasa, bitta oy uchun job ikki marta yurib ketmasligi uchun PostgreSQL advisory lock ishlatilgan:

```text
payout-job:<YYYY-MM>
```

Lockni olgan instance hisoblaydi, boshqalari skip qiladi

## Redis cache

Redis payoutlarni o'qish uchun ishlatiladi.

```text
payout:{id}
TTL: REDIS_PAYOUT_TTL, default 10m
data: payout + adjustments

payouts:courier:{courier_id}:ver
TTL: yo'q
data: kuryer payout list cache versiyasi

payouts:courier:{courier_id}:{ver}:{filters}
TTL: REDIS_PAYOUT_LIST_TTL, default 2m
data: kuryer payout tarixi sahifasi
```

Invalidatsiya:

```text
payout yaratilsa          -> courier list version oshiriladi
payout status o'zgarsa    -> payout key o'chiriladi, list version oshiriladi
adjustment yaratilsa      -> payout key o'chiriladi, list version oshiriladi
```

List cache keylari version orqali ajratilgan. Version o'zgarsa eski list keylar ishlatilmaydi va TTL bilan o'zi tugaydi.

## Errorlar

```text
401  token yo'q yoki noto'g'ri
403  ruxsat yo'q
404  topilmadi
400  noto'g'ri UUID, sana yoki oy formati
400  kelajakdagi oy uchun payout so'ralgan
409  duplicate courier phone
409  duplicate payout calculation
```

## Database uchun ERD

https://dbdiagram.io/d/6aa81964957fec6d5bef586f
