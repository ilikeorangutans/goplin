package database

var migrations []string = []string{
	`
	  create table notebooks (
		id text not null primary key,
		parent_id text not null default "",
		source text not null,
		title text not null,
		sort_order integer not null default 0,
		created_time integer not null,
		updated_time integer not null
	  );
	`,
	`
	  create index idx_folders_title on notebooks (title);
	  create index idx_folders_updated_time on notebooks (updated_time);
	`,
	`
	  create table notes (
		id text not null primary key,
		parent_id text not null default "",
		title text not null default "",
		body text not null default "",
		created_time integer not null,
		updated_time integer not null,
		source text not null,
		todo integer not null default 0
	  );
	`,
	`
	  create table sync_items (
		id string not null primary key
	  );
	`,
}
