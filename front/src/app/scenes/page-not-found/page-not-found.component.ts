import {Component, OnInit} from '@angular/core';
import {MetaService} from '../../services/meta.service';
import {META_TITLES} from '../../utils/constants';

@Component({
  selector: 'app-page-not-found',
  templateUrl: './page-not-found.component.html',
  styleUrls: ['./page-not-found.component.css']
})
export class PageNotFoundComponent implements OnInit {
  constructor(private metaService: MetaService) {
  }

  ngOnInit(): void {
    this.metaService.setTitle(META_TITLES.NOT_FOUND.title);
  }
}
