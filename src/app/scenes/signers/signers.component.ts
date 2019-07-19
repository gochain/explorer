/*CORE*/
import {Component, OnInit} from '@angular/core';
/*SERVICE*/
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {SignerData, SignerStat} from '../../models/signer-stats';
import {ChartItem} from '../../models/chart';
import {sortObjArrByKey} from '../../utils/functions';

interface IHoveredItem {
  rangeIndex: number;
  itemIndex: number;
}

@Component({
  selector: 'app-signers',
  templateUrl: './signers.component.html',
  styleUrls: ['./signers.component.scss']
})
export class SignersComponent implements OnInit {

  statsData: SignerStat[];

  colorScheme = {
    name: 'cool',
    selectable: true,
    group: 'Ordinal',
    domain: [
      '#a8385d', '#7aa3e5', '#a27ea8', '#aae3f5', '#adcded', '#a95963', '#8796c0', '#7ed3ed', '#50abcc', '#ad6886'
    ]
  };

  hoveredItem: IHoveredItem;

  static formChartData(items: SignerData[]): ChartItem[] {
    return items.map((item: SignerData, index: number) => ({
      name: item.name || item.signer,
      value: item.blocks_count,
      extra: {
        itemIndex: index,
      },
    }));
  }

  constructor(
    private _commonService: CommonService,
  ) {
  }

  ngOnInit() {
    this._commonService.getSignerStats().subscribe(data => {
      this.onSignerData(data);
    });
  }

  onSignerData(data: SignerStat[]) {
    this.statsData = data;
    this.processSignerData();
  }

  processSignerData() {
    this.statsData.forEach((stat: SignerStat) => {
      stat.totalBlocks = 0;
      stat.chartData = [];
      stat.signer_stats.forEach((signer: SignerData) => {
        stat.totalBlocks += signer.blocks_count;
        stat.chartData = SignersComponent.formChartData(stat.signer_stats);
      });
      stat.signer_stats.forEach((signer: SignerData) => {
        signer.percent = (signer.blocks_count / stat.totalBlocks * 100).toFixed(4);
      });

      // default sorting by block count
      sortObjArrByKey(stat.signer_stats, 'blocks_count');
    });
  }

  onChartItemEnter(data: any, rangeIndex) {
    this.hoveredItem = {
      rangeIndex,
      itemIndex: data.value.extra.itemIndex,
    };
  }

  onChartItemLeave() {
    this.hoveredItem = null;
  }
}
