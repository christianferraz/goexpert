package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/christianferraz/goexpert/humap/busca_aghu_ortopedia/configs"
	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
)

// Cirurgia representa uma estrutura para dados de cirurgias.
type InformacoesCirurgia struct {
	AgdSeq               sql.NullInt64  `json:"agd_seq"`
	Seq                  sql.NullInt64  `json:"seq"`
	PacCodigo            sql.NullInt64  `json:"pac_codigo"`
	Prontuario           sql.NullInt64  `json:"prontuario"`
	Paciente             sql.NullString `json:"paciente"`
	DtNascimentoPaciente sql.NullString `json:"dt_nascimento_paciente"` // TEXT pode ser nulo em alguns DBs
	Age                  sql.NullString `json:"age"`                    // INTERVAL é um tipo complexo, representado como string
	Unidade              sql.NullString `json:"unidade"`
	Medico               sql.NullString `json:"medico"`
	NomeEspecialidade    sql.NullString `json:"nome_especialidade"`
	EspSeq               sql.NullInt16  `json:"esp_seq"`
	ProcedimentoCirurgia sql.NullString `json:"procedimento_cirurgia"`
	Data                 sql.NullTime   `json:"data"`
	NroAgenda            sql.NullInt16  `json:"nro_agenda"`
	Situacao             sql.NullString `json:"situacao"`
	NaturezaAgend        sql.NullString `json:"natureza_agend"`
	OrigemPacCirg        sql.NullString `json:"origem_pac_cirg"`
	SciSeqp              sql.NullInt16  `json:"sci_seqp"`
	AtdSeq               sql.NullInt64  `json:"atd_seq"`
	DpaSeq               sql.NullInt16  `json:"dpa_seq"`
	DthrPrevInicio       sql.NullTime   `json:"dthr_prev_inicio"`
	DthrPrevFim          sql.NullTime   `json:"dthr_prev_fim"`
	TempoUtlzO2          sql.NullInt16  `json:"tempo_utlz_o2"`
	TempoUtlzProAzot     sql.NullInt16  `json:"tempo_utlz_pro_azot"`
	DthrInicioAnest      sql.NullTime   `json:"dthr_inicio_anest"`
	DthrFimAnest         sql.NullTime   `json:"dthr_fim_anest"`
	DthrInicioCirg       sql.NullTime   `json:"dthr_inicio_cirg"`
	DthrFimCirg          sql.NullTime   `json:"dthr_fim_cirg"`
	DthrEntradaSala      sql.NullTime   `json:"dthr_entrada_sala"`
	DthrSaidaSala        sql.NullTime   `json:"dthr_saida_sala"`
	DthrEntradaSr        sql.NullTime   `json:"dthr_entrada_sr"`
	DthrSaidaSr          sql.NullTime   `json:"dthr_saida_sr"`
	DthrDigitNotaSala    sql.NullTime   `json:"dthr_digit_nota_sala"`
	Asa                  sql.NullInt16  `json:"asa"`
	DthrUltAtlzNotaSala  sql.NullTime   `json:"dthr_ult_atlz_nota_sala"`
	TempoPrevHrs         sql.NullInt16  `json:"tempo_prev_hrs"`
	TempoPrevMin         sql.NullInt16  `json:"tempo_prev_min"`
	CctCodigo            sql.NullInt64  `json:"cct_codigo"`
	IndTemDescricao      sql.NullString `json:"ind_tem_descricao"`
	ComplementoCanc      sql.NullString `json:"complemento_canc"`
	IndOverbooking       sql.NullString `json:"ind_overbooking"`
	DthrInicioOrdem      sql.NullTime   `json:"dthr_inicio_ordem"`
	MomentoAgend         sql.NullString `json:"momento_agend"`
	UtilizacaoSala       sql.NullString `json:"utilizacao_sala"`
	Version              sql.NullInt64  `json:"version"`
	IndAplLstCrgSeg      sql.NullString `json:"ind_apl_lst_crg_seg"`
	IndPrc               sql.NullString `json:"ind_prc"`
	AtbProf              sql.NullString `json:"atb_prof"`
	DthrAtbProf          sql.NullTime   `json:"dthr_atb_prof"`
}

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	host := config.DBHost
	port := config.DBPort
	user := config.DBUser
	password := config.DBPassword
	dbname := config.DBName

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected!")
	GetXls(db, 0)
}

func GetXls(db *sql.DB, prontuario int64) (string, error) {
	arquivo, err := xlsx.OpenFile("cirorto2023faturado_20231228.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	planilha := arquivo.Sheet["cirorto2023faturado_20231228"]
	if planilha == nil {
		log.Fatal("Planilha não encontrada")
		return "", fmt.Errorf("Planilha não encontrada")
	}

	for _, row := range planilha.Rows {
		if row != nil && len(row.Cells) > 0 {
			cell := row.Cells[0] // Célula da coluna A
			text := cell.String()
			nome, _ := GetCNS(db, text)

			row.Cells[1].Value = nome

		}

	}
	arquivo.Save("cirorto2023faturado_202312282.xlsx")

	// err := db.QueryRow(`SELECT cns FROM agh.aip_pacientes WHERE codigo = $1`, prontuario).Scan(&cns)
	// if err != nil {
	// 	return "", err
	// }

	return "cns", nil
}

func GetCNS(db *sql.DB, cns string) (string, error) {
	var pessoa string

	row := db.QueryRow(`SELECT pes.nome FROM agh.rap_pessoa_tipo_informacoes AS tii
	LEFT JOIN agh.rap_pessoas_fisicas pes ON pes.codigo = tii.pes_codigo
	WHERE tii.tii_seq = 7 AND tii.valor = $1`, cns)
	err := row.Scan(&pessoa)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Nenhum resultado encontrado.")
			return "", err
		}
		log.Println(err.Error())
		return "", err
	}
	return pessoa, nil
}

func GetInformacoesCirurgia(db *sql.DB, prontuario int64) ([]InformacoesCirurgia, error) {
	rows, err := db.Query(`select 
		agd.seq AS agd_seq,
		cir.seq,
		cir.pac_codigo,
		pacientes.prontuario, 
		pacientes.nome AS Paciente,
		TO_CHAR(pacientes.dt_nascimento, 'DD/MM/YYYY') AS dt_nascimento_paciente,
		age(pacientes.dt_nascimento),
		uni_func.descricao AS Unidade,
		PF.nome AS Medico,
		esp.nome_especialidade,
		cir.esp_seq,
		pci.descricao AS procedimento_cirurgia, 
		cir.data, 
		cir.nro_agenda, 
		cir.situacao, 
		cir.natureza_agend, 
		cir.origem_pac_cirg, 
		cir.sci_seqp, 
		cir.atd_seq, 
		cir.dpa_seq, 
		cir.dthr_prev_inicio, 
		cir.dthr_prev_fim, 
		cir.tempo_utlz_o2, 
		cir.tempo_utlz_pro_azot, 
		cir.dthr_inicio_anest, 
		cir.dthr_fim_anest, 
		cir.dthr_inicio_cirg, 
		cir.dthr_fim_cirg, 
		cir.dthr_entrada_sala, 
		cir.dthr_saida_sala, 
		cir.dthr_entrada_sr, 
		cir.dthr_saida_sr, 
		cir.dthr_digit_nota_sala, 
		cir.asa, 
		cir.dthr_ult_atlz_nota_sala, 
		cir.tempo_prev_hrs, 
		cir.tempo_prev_min, 
		cir.cct_codigo, 
		cir.ind_tem_descricao, 
		cir.complemento_canc, 
		cir.ind_overbooking, 
		cir.dthr_inicio_ordem, 
		cir.momento_agend, 
		cir.utilizacao_sala, 
		cir.version, 
		cir.ind_apl_lst_crg_seg, 
		cir.ind_prc, 
		cir.atb_prof, 
		cir.dthr_atb_prof
	FROM agh.mbc_cirurgias cir
	INNER JOIN agh.mbc_proc_descricoes proc_desc ON cir.seq = proc_desc.dcg_crg_seq
	INNER JOIN agh.mbc_procedimento_cirurgicos proc_cir ON proc_desc.dcg_seqp = proc_cir.seq 
	INNER JOIN agh.agh_unidades_funcionais uni_func ON uni_func.seq = cir.unf_seq
	INNER JOIN agh.mbc_agendas agd ON agd.seq = cir.agd_seq
	INNER JOIN agh.mbc_proc_esp_por_cirurgias ppc ON ppc.crg_seq = cir.seq
	INNER JOIN agh.mbc_especialidade_proc_cirgs epr ON epr.pci_seq = ppc.epr_pci_seq
	INNER JOIN agh.mbc_procedimento_cirurgicos pci ON pci.seq = epr.pci_seq
	INNER JOIN agh.agh_especialidades esp ON esp.seq = cir.esp_seq
	INNER JOIN agh.rap_servidores servidores ON servidores.matricula = agd.puc_ser_matricula AND servidores.vin_codigo = agd.puc_ser_vin_codigo
	INNER JOIN agh.rap_pessoas_fisicas PF ON servidores.pes_codigo = PF.codigo
	INNER JOIN agh.aip_pacientes pacientes ON pacientes.codigo = cir.pac_codigo
	WHERE cir.dthr_inicio_ordem BETWEEN '2023-01-01' AND '2023-12-31' LIMIT 10`)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	var cirurgias []InformacoesCirurgia

	for rows.Next() {
		var cirurgia InformacoesCirurgia

		err := rows.Scan(
			&cirurgia.AgdSeq,
			&cirurgia.Seq,
			&cirurgia.PacCodigo,
			&cirurgia.Prontuario,
			&cirurgia.Paciente,
			&cirurgia.DtNascimentoPaciente,
			&cirurgia.Age,
			&cirurgia.Unidade,
			&cirurgia.Medico,
			&cirurgia.NomeEspecialidade,
			&cirurgia.EspSeq,
			&cirurgia.ProcedimentoCirurgia,
			&cirurgia.Data,
			&cirurgia.NroAgenda,
			&cirurgia.Situacao,
			&cirurgia.NaturezaAgend,
			&cirurgia.OrigemPacCirg,
			&cirurgia.SciSeqp,
			&cirurgia.AtdSeq,
			&cirurgia.DpaSeq,
			&cirurgia.DthrPrevInicio,
			&cirurgia.DthrPrevFim,
			&cirurgia.TempoUtlzO2,
			&cirurgia.TempoUtlzProAzot,
			&cirurgia.DthrInicioAnest,
			&cirurgia.DthrFimAnest,
			&cirurgia.DthrInicioCirg,
			&cirurgia.DthrFimCirg,
			&cirurgia.DthrEntradaSala,
			&cirurgia.DthrSaidaSala,
			&cirurgia.DthrEntradaSr,
			&cirurgia.DthrSaidaSr,
			&cirurgia.DthrDigitNotaSala,
			&cirurgia.Asa,
			&cirurgia.DthrUltAtlzNotaSala,
			&cirurgia.TempoPrevHrs,
			&cirurgia.TempoPrevMin,
			&cirurgia.CctCodigo,
			&cirurgia.IndTemDescricao,
			&cirurgia.ComplementoCanc,
			&cirurgia.IndOverbooking,
			&cirurgia.DthrInicioOrdem,
			&cirurgia.MomentoAgend,
			&cirurgia.UtilizacaoSala,
			&cirurgia.Version,
			&cirurgia.IndAplLstCrgSeg,
			&cirurgia.IndPrc,
			&cirurgia.AtbProf,
			&cirurgia.DthrAtbProf,
		)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		cirurgias = append(cirurgias, cirurgia)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	for _, c := range cirurgias {
		dd := c.DthrUltAtlzNotaSala.Time.Format("02/01/2006 15:04:05")

		fmt.Printf("%+v\n\n", dd)
	}
	return cirurgias, nil
}
