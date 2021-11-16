

```sql

select a.supplier_type,a.activity_date,count(1) times,
b.duration,
b.guest_quantity,
b.category_loc,b.priced_category_code,b.agent_id,b.region_code,b.config_field,b.ship_code
from ttd_order a,ttd_cruise_mapping b
where a.relation_id = b.id and a.activity_date > '2021-10-22' and a.supplier_type in (1,2)
group by a.supplier_type,a.activity_date,
b.duration,
b.guest_quantity,
b.category_loc,b.priced_category_code,b.agent_id,b.region_code,b.config_field,b.ship_code
order by times desc


```


```sql

select id,name,text_val from klraildb.ptp_global_config where name like 'cruise%'

```