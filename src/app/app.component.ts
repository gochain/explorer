import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ApiService } from './api.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  searchFor: string = '';
  constructor(private api: ApiService,private router: Router) {    
  }
  ngOnInit() {
  
  }
  
  search() {
    console.log(this.searchFor);   
    if (this.searchFor.length == 42){
    this.router.navigate(['/address/',this.searchFor]); 
    }
    if (this.searchFor.length == 66){
      this.router.navigate(['/tx/',this.searchFor]); 
      }
  }
    
}
