package model

var migrations []string = []string{
	`
	create table notebooks (
	  id text not null primary key,
	  parent_id text,
	  source text not null,
	  title text not null,
	  sort_order integer not null default 0
	);
	`,
}
