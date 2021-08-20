package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type Point struct {
	Timestamp string `json:"timestamp"`
	CPU       uint64 `json:"cpu"`
	Memory    uint64 `json:"memory"`
}

/*
	CreateDatabase creates tables for node and pod metrics
*/
func CreateDatabase(db *sql.DB) error {
	sqlStmt := `
	create table if not exists nodes (cluster text, uid text, name text, cpu text, memory text, storage text, time datetime);
	create table if not exists pods (cluster text, uid text, name text, namespace text, container text, cpu text, memory text, storage text, time datetime);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

/*
	UpdateDatabase updates nodeMetrics and podMetrics with scraped data
*/
func UpdateDatabase(db *sql.DB, cluster string, nodeMetrics *v1beta1.NodeMetricsList, podMetrics *v1beta1.PodMetricsList) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into nodes(cluster, uid, name, cpu, memory, storage, time) values(?, ?, ?, ?, ?, ?, datetime('now'))")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range nodeMetrics.Items {
		_, err = stmt.Exec(cluster, v.UID, v.Name, v.Usage.Cpu().MilliValue(), v.Usage.Memory().MilliValue()/1000, v.Usage.StorageEphemeral().MilliValue()/1000)
		if err != nil {
			return err
		}
	}

	stmt, err = tx.Prepare("insert into pods(cluster, uid, name, namespace, container, cpu, memory, storage, time) values(?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))") // customized by kore-board (add table column "cluster")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range podMetrics.Items {
		for _, u := range v.Containers {
			_, err = stmt.Exec(cluster, v.UID, v.Name, v.Namespace, u.Name, u.Usage.Cpu().MilliValue(), u.Usage.Memory().MilliValue()/1000, u.Usage.StorageEphemeral().MilliValue()/1000) // customized by kore-board (add table column "cluster")
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit()

	if err != nil {
		rberr := tx.Rollback()
		if rberr != nil {
			return rberr
		}
		return err
	}

	return nil
}

/*
	CullDatabase deletes rows from nodes and pods based on a time window.
*/
func CullDatabase(db *sql.DB, cluster string, window *time.Duration) error { // customized by kore-board
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	windowStr := fmt.Sprintf("-%.0f seconds", window.Seconds())

	nodestmt, err := tx.Prepare("delete from nodes where cluster=? and time <= datetime('now', ?);") // customized by kore-board
	if err != nil {
		return err
	}

	defer nodestmt.Close()
	res, err := nodestmt.Exec(cluster, windowStr) // customized by kore-board
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	log.Debugf("Cleaning up nodes: %d rows removed", affected)

	podstmt, err := tx.Prepare("delete from pods where  cluster=? and time <= datetime('now', ?);") // customized by kore-board
	if err != nil {
		return err
	}

	defer podstmt.Close()
	res, err = podstmt.Exec(cluster, windowStr) // customized by kore-board
	if err != nil {
		return err
	}

	affected, _ = res.RowsAffected()
	log.Debugf("Cleaning up pods on %s: %d rows removed", cluster, affected) // customized by kore-board
	err = tx.Commit()

	if err != nil {
		rberr := tx.Rollback()
		if rberr != nil {
			return rberr
		}
		return err
	}

	return nil
}

func Select(db *sql.DB, table string, cluster string, namespace string, names string, op string) ([]Point, error) {

	if op == "" {
		op = "SUM"
	}

	sql := fmt.Sprintf("SELECT time, CAST(%s(cpu) AS INTEGER) cpu, CAST(%s(memory) AS INTEGER) memory FROM %s WHERE cluster=?", op, op, table)
	params := []interface{}{cluster}

	if table == "pods" {
		sql = sql + " AND namespace=?"
		params = append(params, namespace)
	}

	// name
	if names != "" {
		// multi
		if strings.ContainsAny(names, ",") {
			subargs := []string{}
			for _, v := range strings.Split(names, ",") {
				subargs = append(subargs, "?")
				params = append(params, v)
			}
			sql = sql + fmt.Sprintf(" AND name IN (%s)", strings.Join(subargs, ","))
		} else {
			// single
			params = append(params, names)
			sql = sql + fmt.Sprintf(" AND name = ?")
		}
	}

	sql = sql + " GROUP BY time ORDER BY time"

	log.Infof("sql=%s, params=%v", sql, params)
	rows, err := db.Query(sql, params...)
	if err != nil {
		log.Errorf("Error getting pod metrics: %v", err)
		return nil, err
	}

	defer rows.Close()

	resultList := []Point{}

	for rows.Next() {
		var cpu string
		var memory string
		var timestamp string
		err = rows.Scan(&timestamp, &cpu, &memory)
		if err != nil {
			return nil, err
		}

		//cpu
		c, err := strconv.ParseUint(cpu, 10, 64)
		if err != nil {
			log.Errorln(err)
			c = 0
		}
		//memory
		m, err := strconv.ParseUint(memory, 10, 64)
		if err != nil {
			log.Errorln(err)
			m = 0
		}

		resultList = append(resultList, Point{
			Timestamp: timestamp,
			CPU:       c,
			Memory:    m,
		})

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return resultList, nil

}
