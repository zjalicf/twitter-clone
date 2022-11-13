import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { MainPageComponent } from './components/main-page/main-page.component';
import { RegisterBusinessComponent } from './components/register-business/register-business.component';
import { RegisterRegularComponent } from './components/register-regular/register-regular.component';

const routes: Routes = [
  {
    path: 'Main-Page',
    component: MainPageComponent
  },
  {
    path: 'Register-Regular',
    component: RegisterRegularComponent
  },
  {
    path: 'Register-Business',
    component: RegisterBusinessComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
