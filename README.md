# Общая информация о проекте
---
> Проект разрабатывается в рамках предмета "Обеспечение качества и тестирования ПО" в Уфимском Гос. Авиационном Техническом Университете
---
Здесь мы разрабатываем программный модуль __«Учет нарушений правил дорожного движения»__. Для каждой автомашины (и ее владельца) в базе хранится список нарушений. Для каждого нарушения фиксируется дата, время, вид нарушения и размер штрафа. При оплате всех штрафов машина удаляется из базы.

### Структура проекта

Весь код будет находится в папке *internal*, внутри которой код с бекенд и фронтенд будет разделен на папки соответственно *back* и *front*.

Весь процесс разработки делится на этапы:
1. Проектирование системы. Разработка предложений по реализации системы;
2. Непосредственно разработка программного модуля «Учет нарушений правил дорожного движения».
3. Тестирование и отладка модуля;
4. Внедрение и сопровождение (на самом деле хз что как и куда). Сопровождение также включает в себя:
4.1. Руководство пользователя;
4.2. Комментарии в коде, там где это необходимо.

### Стек технологий

Весь интерфейс, то, что будет представлено перед пользователями, будет написано на 
- HTML
- CSS

Серверная часть будет написана на
- Golang
- PHP (там, где моих знаний Go будет недостаточно, с работой на *php* есть небольшой опыт)