import argparse
import asyncio
import sys

from rich.console import Console

from fatura.config import load_config
from fatura.exceptions import ConfigError
from fatura.logging_config import setup_logging

console = Console()


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="5G Energia - Automação de Faturas Coelba/Neoenergia",
    )
    parser.add_argument(
        "--config", "-c",
        default="config.yaml",
        help="Caminho para o arquivo de configuração YAML",
    )
    parser.add_argument(
        "--mes-ano", "-m",
        help="Mês/ano no formato MMAAAA (ex: 032026). Se omitido, baixa a fatura mais recente disponível.",
    )
    parser.add_argument(
        "--uc",
        help="Processar apenas esta UC (para testes)",
    )
    parser.add_argument(
        "--force", "-f",
        action="store_true",
        help="Re-baixar mesmo se já existir no banco",
    )
    parser.add_argument(
        "--headed",
        action="store_true",
        help="Rodar browser com interface gráfica (modo debug)",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Apenas validar configuração e sair",
    )
    parser.add_argument(
        "--verbose", "-v",
        action="count",
        default=0,
        help="Aumentar nível de log (-v info, -vv debug)",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    setup_logging(args.verbose)

    try:
        config = load_config(args.config)
    except ConfigError as e:
        console.print(f"[red]Erro de configuração:[/red] {e}")
        sys.exit(1)

    if args.headed:
        config.portal.headless = False

    clientes_ativos = [c for c in config.clientes if c.ativo]
    if args.uc:
        clientes_ativos = [c for c in clientes_ativos if c.uc == args.uc]
        if not clientes_ativos:
            console.print(f"[red]UC {args.uc} não encontrada na configuração[/red]")
            sys.exit(1)

    console.print(f"[green]Configuração carregada:[/green] {len(clientes_ativos)} cliente(s) ativo(s)")

    if args.dry_run:
        console.print("[yellow]Modo dry-run: configuração válida. Saindo.[/yellow]")
        for c in clientes_ativos:
            tipo = "CNPJ" if c.is_cnpj else "CPF"
            console.print(f"  - {c.nome} | UC: {c.uc} | {tipo} | {c.tipo_acesso.value}")
        sys.exit(0)

    from fatura.jobs import processar_faturas_mes

    resultado = asyncio.run(
        processar_faturas_mes(
            config=config,
            clientes=clientes_ativos,
            mes_ano=args.mes_ano,
            force=args.force,
        )
    )

    console.print()
    console.print("[bold]Resultado:[/bold]")
    console.print(f"  Total:   {resultado.total}")
    console.print(f"  Sucesso: [green]{resultado.sucesso}[/green]")
    console.print(f"  Erro:    [red]{resultado.erro}[/red]")
    console.print(f"  Pulado:  [yellow]{resultado.pulado}[/yellow]")

    if resultado.erro > 0:
        sys.exit(1)


if __name__ == "__main__":
    main()
