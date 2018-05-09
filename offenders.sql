select
	case
		when business = '' and owner = '' then '*unknown*'
		when business = '' then owner
		when owner = '' then business
		else business || '/' || owner
	end as businessowner,
	sum(amount)
from
	taxes
group by
	businessowner
order by
	sum desc;
