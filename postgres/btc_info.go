package postgres

import (
	"fmt"
	"time"
)

func (store *Storage) InsertBTCInfo(btcInfo BTCInfo) (BTCInfo, error) {
	data := BTCInfo{}
	stmt, err := store.db.PrepareNamed(`
		INSERT INTO btc_info (
			amount,
			created,
			offsettz
		) VALUES (
			:amount,
			:created,
			:offsettz
		) RETURNING *
    `)
	defer func() {
		if stmt != nil {
			if err := stmt.Close(); err != nil {
				return
			}
		} else {
			return
		}
	}()
	if err != nil {
		return data, err
	}
	err = stmt.Get(&data, btcInfo)
	return data, err
}

func (store *Storage) GetBTCHistory(startDate time.Time, endDate time.Time) ([]BTCInfoResult, error) {
	data := make([]BTCInfoResult, 0)

	query := fmt.Sprintf(`with dates as (
			select
				generate_series(
					(date '%v') :: timestamptz,
					(date '%v') :: timestamptz,
					interval '1 hour'
				) as dt
		)
		select
			(
				d.dt :: date :: text || ' ' || to_char(d.dt :: time, 'HH24:MM:SS')
			) :: timestamptz as created,
			coalesce(sum(b.amount), 0) as amount
		from
			dates d
			left join btc_info b on b.created >= d.dt
			and b.created <= d.dt + interval '1 hour'
		group by
			d.dt
		having
			sum(b.amount) > 0
		order by
			d.dt;`, startDate.Format(dateFmt), endDate.Format(dateFmt))
	if err := store.db.Select(&data, query); err != nil {
		return nil, err
	}
	return data, nil
}
