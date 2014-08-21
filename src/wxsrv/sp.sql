delimiter $$

create procedure report_all()
	begin
	select sum(er.execisetime) as all_execise_time, sum(er.execiseenergy) as all_execise_energy from execise_records as er;
	end $$

create procedure report_this_week()
	begin
	select sum(er.execisetime) as all_execise_time, sum(er.execiseenergy) as 
		all_execise_energy from execise_records as er where week(er.record_time)=week(current_time());
	end $$

delimiter ;